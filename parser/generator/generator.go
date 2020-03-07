// Command generator creates machine.go from machine.txt
package main

import (
	"bufio"
	"bytes"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"sethwklein.net/go/errors"
)

func mainError() (err error) {
	txt, err := ioutil.ReadFile("machine.txt")
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(bytes.NewReader(txt))
	defer func() {
		// can't see how there would be any
		err = errors.Append(err, scanner.Err())
	}()

	buffer := &bytes.Buffer{}

	buffer.WriteString(`package parser

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

	`)

	first := true
	open := false
	openSwitch := func(label []byte) {
		if first {
			buffer.WriteString("// ")
		}
		buffer.Write(label)
		buffer.WriteString(":\n")
		if first {
			buffer.WriteString("i := 0\n")
			first = false
		} else {
			buffer.WriteString("i++\n")
		}
		buffer.WriteString("switch tokens[i] {\n")
		open = true
	}
	closeSwitch := func() {
		buffer.WriteString(`default:
				panic(TokenError{tokens[i], i})
			}
		`)
		open = false
	}

	rCase := regexp.MustCompile(`^([A-Za-z]+)\+([a-z]): ([A-Za-z-]+) +\((.*)\)$`)
	for scanner.Scan() {
		matches := rCase.FindSubmatch(scanner.Bytes())
		if len(matches) < 1 {
			continue
		}
		if !open {
			openSwitch(matches[1])
		}

		buffer.WriteString("case '")
		buffer.Write(matches[2])
		buffer.WriteString("':\n")

		if actions := matches[4]; len(actions) > 0 {
			actions = bytes.Replace(actions, []byte{','},
				[]byte{';'}, -1)
			actions = bytes.Replace(actions, []byte("new page"),
				[]byte("newPage()"), -1)
			actions = bytes.Replace(actions, []byte(".low"),
				[]byte(".Low"), -1)
			actions = bytes.Replace(actions, []byte(".high"),
				[]byte(".High"), -1)
			actions = bytes.Replace(actions, []byte("terminate"),
				[]byte("return pages"), -1)
			buffer.Write(actions)
			buffer.WriteByte('\n')
		}

		if len(matches[3]) > 0 && matches[3][0] != '-' {
			buffer.WriteString("goto ")
			buffer.Write(matches[3])
			buffer.WriteByte('\n')
		} else {
			closeSwitch()
		}
	}

	buffer.WriteString("}")

	// used for debugging
	// os.Stdout.Write(buffer.Bytes())

	formatted, err := format.Source(buffer.Bytes())
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("machine.go", formatted, 0666)
	if err != nil {
		return err
	}
	return nil
}

func mainCode() int {
	err := mainError()
	if err == nil {
		return 0
	}
	fmt.Fprintf(os.Stderr, "%v: Error: %v\n", filepath.Base(os.Args[0]), err)
	return 1
}

func main() {
	os.Exit(mainCode())
}
