package parser

// AUTOMATICALLY GENERATED! DO NOT EDIT!
// recreate with: go run generator/generator.go

import (
	"errors"

	"sethwklein.net/thefile/tokenizer"
)

// Parse finds the pages in tokens.
func Parse(tokens []tokenizer.Token) (pages []Page) {
	if len(tokens) < 1 || tokens[len(tokens)-1] != tokenizer.Default.E {
		panic(errors.New("missing end of file token"))
	}

	var head, body *Part
	newPage := func() {
		pages = append(pages, Page{})
		head = &pages[len(pages)-1].Head
		body = &pages[len(pages)-1].Body
	}
	var note int

	// I:
	i := 0
	switch tokens[i] {
	case 'c':
		goto Ic
	case 't':
		newPage()
		head.Low = i
		goto Q
	case 'o':
		newPage()
		head.empty(i)
		body.Low = i
		goto Bo
	case 'e':
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Ic:
	i++
	switch tokens[i] {
	case 'c':
		note = i
		goto Icc
	case 't':
		newPage()
		head.Low = i
		goto Q
	case 'o':
		newPage()
		head.empty(i)
		body.Low = i
		goto Bo
	case 'e':
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Icc:
	i++
	switch tokens[i] {
	case 'c':
		goto Icc
	case 't':
		newPage()
		head.Low = i
		goto Q
	case 'o':
		newPage()
		head.empty(note)
		body.Low = note
		goto Bo
	case 'e':
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Q:
	i++
	switch tokens[i] {
	case 'c':
		head.High = i
		goto Hc
	case 't':
		goto Q
	case 'o':
		head.High = head.Low
		body.Low = head.Low
		goto Bo
	case 'e':
		head.High = i
		body.empty(i)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Hc:
	i++
	switch tokens[i] {
	case 'c':
		body.Low = i
		goto Hcc
	case 't':
		note = i
		goto Hct
	case 'o':
		body.Low = i
		goto Bo
	case 'e':
		body.empty(head.High)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Hct:
	i++
	switch tokens[i] {
	case 'c':
		body.empty(head.High)
		newPage()
		head.Low = note
		head.High = i
		goto Hc
	case 't':
		goto Hct
	case 'o':
		body.Low = note
		goto Bo
	case 'e':
		body.empty(head.High)
		newPage()
		head.Low = note
		head.High = i
		body.empty(i)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Hcc:
	i++
	switch tokens[i] {
	case 'c':
		goto Bc
	case 't':
		note = i
		goto Hcct
	case 'o':
		goto Bo
	case 'e':
		body.empty(head.High)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Hcct:
	i++
	switch tokens[i] {
	case 'c':
		body.empty(head.High)
		newPage()
		head.Low = note
		head.High = i
		goto Hc
	case 't':
		goto Hcct
	case 'o':
		goto Bo
	case 'e':
		body.empty(head.High)
		newPage()
		head.Low = note
		head.High = i
		body.empty(i)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Bo:
	i++
	switch tokens[i] {
	case 'c':
		goto Bc
	case 't':
		goto Bo
	case 'o':
		goto Bo
	case 'e':
		body.High = i
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Bc:
	i++
	switch tokens[i] {
	case 'c':
		goto Bc
	case 't':
		note = i
		goto Bct
	case 'o':
		goto Bo
	case 'e':
		body.High = i - 1
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
Bct:
	i++
	switch tokens[i] {
	case 'c':
		body.High = note - 1
		newPage()
		head.Low = note
		head.High = i
		goto Hc
	case 't':
		goto Bct
	case 'o':
		goto Bo
	case 'e':
		body.High = note - 1
		newPage()
		head.Low = note
		head.High = i
		body.empty(i)
		return pages
	default:
		panic(TokenError{tokens[i], i})
	}
}
