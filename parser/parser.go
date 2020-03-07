// Package parser finds the pages in a slice of tokens from tokenizer.
//
// A normal page is /t+(co+)*c/, although technically the last clear is not
// part of the page. There are many corner cases, but those are documented
// only in the tests.
//
// This package expects 'a' to be tokenized as 't', in contrast to what the
// tokenizer package produces by default.
package parser

import (
	"fmt"

	"sethwklein.net/thefile/tokenizer"
)

// Part represents a header or body. It contains the values that would slice
// the part's lines out of a slice of all lines in the file.
type Part struct {
	Low, High int
}

func (p *Part) empty(i int) {
	p.Low = i
	p.High = i
}

// Page contains the parts for a page. In the case of empty headers or bodies,
// the empty part is set to zero length (low == high), with the value chosen
// so head.low:body.high slices all lines inside the page.
type Page struct {
	Head, Body Part
}

// TokenError represents an invalid token given for a certain line. It is
// used by Parse when panicking.
type TokenError struct {
	Token tokenizer.Token
	Index int
}

// Error satisfies the error interface.
func (err TokenError) Error() string {
	return fmt.Sprintf("bad token, %c, for line %d", err.Token, err.Index+1)
}
