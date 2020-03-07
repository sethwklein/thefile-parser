package thefile

import (
	"fmt"
)

type Index struct {
	pages   []*Page
	address map[int]*Page
	named   map[string][]*Page
	tagged  map[string][]*Page
	in      map[string][]*Page
}

// Pages returns the pages used to create index.
func (index *Index) Pages() []*Page {
	return index.pages
}

// Address returns the page with the given address or nil.
func (index *Index) Address(address int) *Page {
	return index.address[address]
}

type NotFoundError struct {
	Name string
}

func (err NotFoundError) Error() string {
	return "page not found named: " + err.Name
}

type AmbiguousNameError struct {
	Name  string
	Pages []*Page
}

func (err AmbiguousNameError) Error() string {
	return fmt.Sprintf("found %d pages named: %s", len(err.Pages), err.Name)
}

// BUG: move index.Named example into tests
// BUG: devise tests to ensure that index.Named never returns another error
/*
	page, err := index.Named("example")
	switch err.(type) {
	case NotFoundError:
		fmt.Println(err.Name)
	case AmbiguousNameError:
		fmt.Println(err.Name, err.Pages)
	case nil:
	}
*/

// Named returns the first page with name in the first title line or nil if
// there were none. It also returns an error if there was not exactly one. If
// there were zero, the error will be a NotFoundError. If there was more than
// one, it will be an AmbiguousNameError containing all the pages. Named
// returns no other type of error.
func (index *Index) Named(name string) (*Page, error) {
	pages := index.named[name]
	switch len(pages) {
	case 0:
		return nil, NotFoundError{name}
	case 1:
		return pages[0], nil
	default:
		return pages[0], AmbiguousNameError{name, pages}
	}
}

// AllNamed is like Named, but returns all the things.
func (index *Index) AllNamed(name string) []*Page {
	return index.named[name]
}

// Tagged returns all pages with tag at any title position.
func (index *Index) Tagged(tag string) []*Page {
	return index.tagged[tag]
}

// Tags return all the titles in the file. Because it's so rarely used, it
// allocates and fills an array. Cache the result instead of calling it
// repeatedly.
func (index *Index) Tags() []string {
	tags := make([]string, 0, len(index.tagged))
	for tag, _ := range index.tagged {
		tags = append(tags, tag)
	}
	return tags
}

// In returns all pages with the given thing at any title position but
// the first. Titles that duplicate the first are not returned.
func (index *Index) In(thing string) []*Page {
	return index.in[thing]
}

func NewIndex(pages []*Page) *Index {
	address := make(map[int]*Page)
	named := make(map[string][]*Page)
	tagged := make(map[string][]*Page)
	in := make(map[string][]*Page)
	for _, page := range pages {
		name, anonymous := page.Name()
		if !anonymous {
			named[name] = append(named[name], page)
		}

		address[page.Address()] = page

		for _, tag := range page.Tags() {
			tagged[tag] = append(tagged[tag], page)
		}
		for _, thing := range page.In() {
			in[thing] = append(in[thing], page)
		}
	}
	return &Index{
		pages:   pages,
		address: address,
		named:   named,
		tagged:  tagged,
		in:      in,
	}
}
