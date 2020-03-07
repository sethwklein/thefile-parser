// Package thefile gives structured access to the file.
package thefile

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"unicode"

	"sethwklein.net/thefile/parser"
	"sethwklein.net/thefile/storage"
	"sethwklein.net/thefile/tokenizer"
)

// Page provides access to a page in the file.
type Page struct {
	// titles contains all titles with duplicates removed (the first is kept)
	// and zero length titles removed unless they are in the first position.
	// It always contains at least one title.
	titles []string

	// file is the original file contents, used for creating body slices.
	file []byte

	// offsets contains the offsets of the body lines, including one for
	// just past the end of the last line.
	offsets []int

	// all contains the offsets of all lines used in constructing the page,
	// including one for just past the end of the last line. It was first
	// added to support a less abstract hash.
	all []int

	// line is the number of first line of the head. Used for Address.
	line int

	// index is the page index. Used for Index.
	index int
}

// I frequently want the first title in string form. I must not forget to
// consider whether to handle anonymous names differently.

// Name returns the first title and true if the page is anonymous (if its name
// is zero length).
func (page *Page) Name() (string, bool) {
	name := page.titles[0]
	return name, len(name) < 1
}

// I frequently want to for-range over all titles in string form, including the
// first, but skipping zero length titles and duplicate titles.

// Tags returns all titles not zero length. Duplicated titles appear only
// the first time.
func (page *Page) Tags() []string {
	tags := page.titles
	if len(tags[0]) < 1 {
		tags = tags[1:]
	}
	return tags
}

// Tagged returns whether page is tagged with q.
func (page *Page) Tagged(q string) bool {
	if q == "" {
		return false
	}
	for _, tag := range page.Tags() {
		if tag == q {
			return true
		}
	}
	return false
}

// I frequently want to for-range over all titles, except the first, in string
// form, skipping duplicate titles, but don't care if they're unique in the
// file. Duplicate removal should consider the first title, even though it is
// not returned.

/*
	Naming is hard:

	? in: short! in what? but the answer is, in whatever you want to call
		the thing it's in. In practice, the argument usually answers the
		question to the satisfaction of the reader's mind
	~ categories: long, but i think this expresses the usage
	- links: what about the unique ones?
	- others: other whats?
	- secondaries: long, secondary whats?
	- alternates: long, confusing because alternate titles used to indicate
		(and might still indicate) uniqueness, alternate whats?
	- alts: inconsistently abbreviated
	- following: long, following whats?
	- titles: confusing because you have to remember that the first title
		line doesn't count
*/

// In returns all titles except the first or any with zero length. Duplicated
// titles appear only the first time, and not at all if they duplicate the
// first title.
func (page *Page) In() []string {
	return page.titles[1:]
}

// Some parsers want the entire body.

// Body returns the body all at once.
func (page *Page) Body() []byte {
	end := page.offsets[len(page.offsets)-1]
	return page.file[page.offsets[0]:end:end]
}

// All returns all bytes used to construct the page. It was first added to
// support a less abstract hash.
func (page *Page) All() []byte {
	end := page.all[len(page.all)-1]
	return page.file[page.all[0]:end:end]
}

// Some parsers (which should probably be replaced) want lines.

// Lines returns the body in lines. Lines include the terminating newline,
// except if the last character in the file isn't a newline. Then that line
// will not be newline terminated. However, it will also not be zero length.
func (page *Page) Lines() [][]byte {
	body := make([][]byte, len(page.offsets)-1)
	for i := range body {
		end := page.offsets[i+1]
		body[i] = page.file[page.offsets[i]:end:end]
	}
	return body
}

// I frequently want the page's position in the file. I want it at line
// resolution in integer format.

/*

	I use it for ordering, for map keys, and for editor navigation.

	It should be numeric (integer) not string because I calculate with it
	more often than I convert it to a string.

	Line (not page) resolution is best. Ordering works with any resolution
	that provides unique, comparable values. The same is true for map keys,
	but editor navigation pretty much requires line resolution.

	Line numbers (one based, as opposed to zero based line indexes) are
	best. Map keys and ordering work with either indexes or numbers. Editor
	navigation, however, requires numbers.

*/

// Address returns the line number (one based) of the page's first title line.
func (page *Page) Address() int {
	return page.line
}

// I want the difference between page indexes for sorting bin metrics.

// Index returns the page index.
func (page *Page) Index() int {
	return page.index
}

// HashRaw returns the bytes of a hash of the page titles and body. This may
// be more stable than the address, making it more useful for the web interface.
func (page *Page) HashRaw() []byte {
	// Technically, this uses no internal knowledge and so doesn't need to
	// be a method of Page.

	// It should be higher performance to use the algorithm in Page.Body
	// to get all the raw bytes and hash that, but that includes the title
	// line indicators which are technically encoding, not content, and
	// should be able to change without changing the hash.

	// No, the title line indicators separate tags, preventing collisions
	// when the suffix of one tag is moved to the prefix of the next.

	// The cost of making memoization thread safe using sync.Once is far
	// higher than just calculating the hash again.

	hash := sha256.New()
	for _, title := range page.Tags() {
		hash.Write([]byte(title))
	}
	hash.Write(page.Body())
	return hash.Sum(nil)
}

// Hash64 returns the base 64 encoding (base64.RawURLEncoding) of the page's
// hash.
func (page *Page) Hash64() string {
	return base64.RawURLEncoding.EncodeToString(page.HashRaw())
}

// parseTitle returns "title" from "----title\n", trimming trailing whitespace.
func parseTitle(line []byte) []byte {
	if len(line) <= 5 {
		return nil
	}
	return bytes.TrimRightFunc(line[4:len(line)-1], unicode.IsSpace)
}

// makeTitles returns what should go in Page.titles.
func makeTitles(buf []byte, offsets []int, head parser.Part) []string {
	length := head.High - head.Low
	// the file can start with body content, creating a page with no header
	if length < 1 {
		return []string{""}
	}

	titles := make([]string, 0, length)

line:
	for i := head.Low; i < head.High; i++ {
		title := string(parseTitle(buf[offsets[i]:offsets[i+1]]))
		if len(title) < 1 && i > head.Low {
			continue
		}
		for _, existing := range titles {
			if title == existing {
				continue line
			}
		}
		titles = append(titles, title)
	}

	return titles
}

func pagesFrom(buf []byte) (pages []*Page, nLines int) {
	// magic constants determined by looking at output of average/average.go.
	// lowering length provides no gains distinguishable from the noise.
	skip := 0
	length := 10000
	if len(buf) < length {
		length = len(buf)
	} else {
		// sample from the middle, hoping to avoid responding too
		// quickly to pattern changes.
		skip = (len(buf) - length) / 2
	}
	count := 0
	for _, c := range buf[skip : skip+length] {
		if c == '\n' {
			count++
		}
	}
	average := length
	if count > 0 {
		average /= count
	}
	estimate := 0
	const fudge = 2
	if average > fudge {
		estimate = len(buf) / (average - fudge)
	}

	tokens := make([]tokenizer.Token, 0, estimate)
	offsets := make([]int, 0, estimate)
	tok := tokenizer.Default
	tok.A = 't'
	for offset := 0; ; {
		token, length := tok.Line(buf[offset:])
		tokens = append(tokens, token)
		offsets = append(offsets, offset)
		if token == tok.E {
			break
		}
		offset += length
	}
	//if cap(tokens) != estimate {
	//	fmt.Println("reallocated")
	//}
	parsed := parser.Parse(tokens)

	backing := make([]Page, len(parsed))
	offsets = append(offsets, len(buf))
	for i, p := range parsed {
		backing[i].titles = makeTitles(buf, offsets, p.Head)
		backing[i].file = buf
		backing[i].offsets = offsets[p.Body.Low : p.Body.High+1]
		backing[i].all = offsets[p.Head.Low : p.Body.High+1]
		backing[i].line = p.Head.Low + 1
		backing[i].index = i
	}
	pages = make([]*Page, len(backing))
	for i := range backing {
		pages[i] = &backing[i]
	}

	return pages, len(tokens)
}

func pages() (pages []*Page, nLines int, err error) {
	buf, err := storage.Load()
	if err != nil {
		return nil, 0, err
	}
	pages, nLines = pagesFrom(buf)
	return pages, nLines, nil
}

// Pages returns the pages.
func Pages() ([]*Page, error) {
	pages, _, err := pages()
	return pages, err
}

// Statistics contains information not available by inspecting the pages.
type Statistics struct {
	// LineCount is the number of lines in the file. Not all of those lines
	// are necessarily part of a page.
	LineCount int
}

// PagesStatistics returns the pages and statistics.
func PagesStatistics() ([]*Page, *Statistics, error) {
	pages, nLines, err := pages()
	return pages, &Statistics{
		LineCount: nLines,
	}, err
}
