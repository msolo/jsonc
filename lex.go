// stream: ( string | raw | comment )*
// comment: line-comment range-comment
// string: '"' ( [^"] | '\"' )* '"'
// raw: [^"]
// non-newline: [^\n]
// line-comment: '//' non-newline
// range-comment: '/*' any '*/'

package jsonc

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch {
	case i.typ == itemEOF:
		return "EOF"
	case i.typ == itemError:
		return "ERR:" + i.val
	case len(i.val) > 10:
		return fmt.Sprintf("%.10q...", i.val)
	}
	return fmt.Sprintf("%q", i.val)
}

const (
	itemError itemType = iota // error occurred;
	// value is text of error
	itemEOF
	itemString
	itemRaw
	itemComment
)

type stateFn func(l *lexer) stateFn

type lexer struct {
	name  string    // used only for error reports.
	input string    // the string being scanned.
	start int       // start position of this item.
	pos   int       // current position in the input.
	width int       // width of last rune read from input.
	items chan item // channel of scanned items.
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:  name,
		input: input,
		items: make(chan item),
	}
	go l.run() // Concurrently run state machine.
	return l, l.items
}

// run lexes the input by executing state functions until
// the state is nil.
func (l *lexer) run() {
	for state := lexStream; state != nil; {
		state = state(l)
	}
	close(l.items) // No more tokens will be delivered.
}

// emit passes an item back to the client.
func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

const eof = -1

// next returns the next rune in the input.
func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width =
		utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

// backup steps back one rune.
// Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// peek returns but does not consume
// the next rune in the input.
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// accept consumes the next rune
// if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set.
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// error returns an error token and terminates the scan
// by passing back a nil pointer that will be the next
// state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.items <- item{
		itemError,
		fmt.Sprintf(format, args...),
	}
	return nil
}

func lexStream(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], "\"") {
			if l.pos > l.start {
				l.emit(itemRaw)
			}
			return lexString // Next state.
		}
		if strings.HasPrefix(l.input[l.pos:], "/") {
			if l.pos > l.start {
				l.emit(itemRaw)
			}
			return lexComment // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemRaw)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexString(l *lexer) stateFn {
	// swallow leading "
	l.next()
	for {
		if strings.HasPrefix(l.input[l.pos:], "\"") {
			l.next() // swallow ending "
			l.emit(itemString)
			return lexStream // Next state.
		}
		// look for escaped \""
		if strings.HasPrefix(l.input[l.pos:], "\\\"") {
			l.next()
			l.next()
			continue
		}
		if l.next() == eof {
			return l.errorf("eof during string parse")
		}
	}
}

func lexComment(l *lexer) stateFn {
	if strings.HasPrefix(l.input[l.pos:], "//") {
		return lexLineComment
	}
	if strings.HasPrefix(l.input[l.pos:], "/*") {
		return lexRangeComment
	}
	return lexStream
}

func lexLineComment(l *lexer) stateFn {
	// swallow //
	l.next()
	l.next()
	for {
		if strings.HasPrefix(l.input[l.pos:], "\n") {
			// don't include trailng \n
			l.emit(itemComment)
			return lexStream // Next state.
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(itemComment)
	}
	l.emit(itemEOF) // Useful to make EOF a token.
	return nil      // Stop the run loop.
}

func lexRangeComment(l *lexer) stateFn {
	// swallow /*
	l.next()
	l.next()
	for {
		if strings.HasPrefix(l.input[l.pos:], "*/") {
			// swallow */
			l.next()
			l.next()
			l.emit(itemComment)
			return lexStream
		}
		if l.next() == eof {
			return l.errorf("unexpected EOF scanning comment")
		}
	}
}
