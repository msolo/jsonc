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

func checkItem(t *testing.T, i item, val string) {
	if i.val != val {
		t.Errorf("expected %s: got %s", val, i.val)
	}
}

func Test(t *testing.T) {
	_, items := lex("test", "{}")
	tl := toSlice(t, items)
	checkItem(t, tl[0], `{}`)
}

func TestLineComment(t *testing.T) {
	_, items := lex("test", "{//}\n")
	tl := toSlice(t, items)
	checkItem(t, tl[0], `{`)
}

func TestRangeComment(t *testing.T) {
	_, items := lex("test", "{/**/}")
	tl := toSlice(t, items)
	checkItem(t, tl[1], `/**/`)
}

func TestRangeCommentInString(t *testing.T) {
	_, items := lex("test", `{"/**/"}`)
	tl := toSlice(t, items)
	checkItem(t, tl[1], `"/**/"`)
}

func TestNestedQuoteInString(t *testing.T) {
	_, items := lex("test", `{"\""}`)
	tl := toSlice(t, items)
	checkItem(t, tl[1], `"\""`)
}

func TestNoCommentTerminator(t *testing.T) {
	_, items := lex("test", "{/*}")
	tl := toSlice(t, items)
	if tl[len(tl)-1].typ != itemError {
		t.Error("expected a parsing error - no comment terminator")
	}
}
