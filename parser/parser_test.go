package parser

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"sethwklein.net/thefile/tokenizer"
)

var tests = []struct {
	input    string
	expected string
}{
	// an empty file contains no sections
	{"e", ""},

	// a file with only whitespace contains no sections
	{"c e", ""},
	{"c c e", ""},
	{"c c c e", ""},

	// a file with header lines, possibly missing leading and trailing
	// clears, contains one section
	{"t e", "(0:1 1:1)"},
	{"t t e", "(0:2 2:2)"},
	{"c t t e", "(1:3 3:3)"},
	{"c t t c e", "(1:3 3:3)"},
	{"c t t c c e", "(1:3 3:3)"},
	{"c t t c c c e", "(1:3 4:5)"},
	{"c t t c o c e", "(1:3 4:5)"},

	// a file with junk but no headers contains one section
	// only skip one clear at beginning of file
	{"o e", "(0:0 0:1)"},
	{"c o e", "(1:1 1:2)"},
	{"c o c e", "(1:1 1:2)"},
	{"c o o c e", "(1:1 1:3)"},
	{"c c o c e", "(1:1 1:3)"},
	{"t t o e", "(0:0 0:3)"},

	// a file with junk before a valid header contains two sections
	{"o c t c e", "(0:0 0:1) (2:3 3:3)"},
	{"c c t e", "(2:3 3:3)"},
	{"c c c t e", "(3:4 4:4)"},
	{"c c c c t e", "(4:5 5:5)"},

	// one and only one clear from top and bottom of body is left out
	{"t c c c e", "(0:1 2:3)"},
	{"t c o c e", "(0:1 2:3)"},
	{"t c c c c e", "(0:1 2:4)"},
	{"t c o o e", "(0:1 2:4)"},
	{"t c o c c e", "(0:1 2:4)"},

	// clear before title at beginning is left out
	{"c t c o c e", "(1:2 3:4)"},
	{"c c t c o c e", "(2:3 4:5)"},

	// two headers makes two sections
	{"t c t e", "(0:1 1:1) (2:3 3:3)"},
	{"t c c t e", "(0:1 1:1) (3:4 4:4)"},
	{"t c c t c e", "(0:1 1:1) (3:4 4:4)"},
	{"t c t c e", "(0:1 1:1) (2:3 3:3)"},

	// mixes of titles and other are bodies
	{"t c t t o e", "(0:1 2:5)"},
	{"t c c o t e", "(0:1 2:5)"},
	{"t c c t t o c o c e", "(0:1 2:8)"},

	// various combinations of title lines after a body work correctly
	{"o c t t o e", "(0:0 0:5)"},
	{"o c t t e", "(0:0 0:1) (2:4 4:4)"},
	{"o c t c e", "(0:0 0:1) (2:3 3:3)"},

	// i suppose we should test sanity!
	//0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6
	{"t t c o o c o o c t c t c o o c e", "(0:2 3:8) (9:10 10:10) (11:12 13:15)"},
}

func TestParse(t *testing.T) {
	for _, test := range tests {
		var input []tokenizer.Token
		for _, t := range test.input {
			if t == ' ' {
				continue
			}
			input = append(input, tokenizer.Token(t))
		}
		pages := Parse(input)
		buf := &bytes.Buffer{}
		for i, page := range pages {
			if i > 0 {
				buf.WriteByte(' ')
			}
			fmt.Fprintf(buf, "(%d:%d %d:%d)", page.Head.Low,
				page.Head.High, page.Body.Low, page.Body.High)
		}
		actual := buf.String()
		if reflect.DeepEqual(test.expected, actual) {
			continue
		}
		t.Errorf("\ninput:    %v\nexpected: %v\nactual:   %v",
			test.input, test.expected, actual)
	}
}

func tokensToString(tokens []tokenizer.Token) string {
	var buf []byte
	buf = append(buf, '{')
	for _, t := range tokens {
		buf = append(buf, byte(t))
	}
	buf = append(buf, '}')
	return string(buf)
}

func testEOF(t *testing.T, tokens []tokenizer.Token) {
	defer func() {
		e := recover()
		if err, ok := e.(error); ok && err.Error() == "missing end of file token" {
			return
		}
		t.Errorf("expected panic from tokens: %v", tokensToString(tokens))
	}()
	Parse(tokens)
}

func TestNone(t *testing.T) {
	testEOF(t, []tokenizer.Token{})
}

func TestMissingEOF(t *testing.T) {
	testEOF(t, []tokenizer.Token{'c'})
}

func badToken(t *testing.T, tokens []tokenizer.Token) {
	defer func() {
		e := recover()
		if e == nil {
			t.Errorf("expected panic from tokens: %v", tokensToString(tokens))
			return
		}
		if te, ok := e.(TokenError); ok {
			if te.Token == tokens[len(tokens)-2] && te.Index == len(tokens)-2 {
				return
			}
			t.Errorf("wrong values in token error\nexpected: %v, %v\nactual: %v, %v", tokens[len(tokens)-2], len(tokens)-2, te.Token, te.Index)
			return
		}
		t.Errorf("wrong bad token error: %T: %v", e, e)
	}()
	Parse(tokens)
}

func TestBadTokens(t *testing.T) {
	// I
	badToken(t, []tokenizer.Token{'x', 'e'})
	// Ic
	badToken(t, []tokenizer.Token{'c', 'x', 'e'})
	// Icc
	badToken(t, []tokenizer.Token{'c', 'c', 'x', 'e'})
	// Q
	badToken(t, []tokenizer.Token{'t', 'x', 'e'})
	// Hc
	badToken(t, []tokenizer.Token{'t', 'c', 'x', 'e'})
	// Hct
	badToken(t, []tokenizer.Token{'t', 'c', 't', 'x', 'e'})
	// Hcc
	badToken(t, []tokenizer.Token{'t', 'c', 'c', 'x', 'e'})
	// Hcct
	badToken(t, []tokenizer.Token{'t', 'c', 'c', 't', 'x', 'e'})
	// Bo
	badToken(t, []tokenizer.Token{'o', 'x', 'e'})
	// Bc
	badToken(t, []tokenizer.Token{'o', 'c', 'x', 'e'})
	// Bct
	badToken(t, []tokenizer.Token{'o', 'c', 't', 'x', 'e'})
}

func TestTokenError(t *testing.T) {
	err := TokenError{'x', 0}
	actual := err.Error()
	expected := "bad token, x, for line 1"
	if actual != expected {
		t.Errorf("unexpected TokenError string\nexpected %v\nactual: %v", expected, actual)
	}
}
