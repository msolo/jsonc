package jsonc

import "testing"

func toSlice(t *testing.T, items chan item) []item {
	tl := make([]item, 0, 16)
	for i := range items {
		tl = append(tl, i)
	}
	t.Logf("tokens: %v", tl)
	return tl
}

func lexToSlice(t *testing.T, s string) []item {
	l := lex("test-lex", s)
	items := make([]item, 0, 16)
	for {
		i := l.yield()
		items = append(items, i)
		if i.typ == itemEOF || i.typ == itemError {
			return items
		}
	}
}

func checkItem(t *testing.T, i item, val string) {
	if i.val != val {
		t.Errorf("expected %s: got %s", val, i.val)
	}
}

func Test(t *testing.T) {
	tl := lexToSlice(t, `{}`)
	checkItem(t, tl[0], `{}`)
}

func TestLineComment(t *testing.T) {
	tl := lexToSlice(t, `{//}\n`)
	checkItem(t, tl[0], `{`)
}

func TestRangeComment(t *testing.T) {
	tl := lexToSlice(t, `{/**/}`)
	checkItem(t, tl[1], `/**/`)
}

func TestRangeCommentInString(t *testing.T) {
	tl := lexToSlice(t, `{"/**/"}`)
	checkItem(t, tl[1], `"/**/"`)
}

func TestNestedQuoteInString(t *testing.T) {
	tl := lexToSlice(t, `{"\""}`)
	checkItem(t, tl[1], `"\""`)
}

func TestNoCommentTerminator(t *testing.T) {
	tl := lexToSlice(t, `{/*}`)
	if tl[len(tl)-1].typ != itemError {
		t.Error("expected a parsing error - no comment terminator")
	}
}
