package thefile

// Intersect returns a list of pages in both a and b. Its output is undefined
// unless both are sorted by address. Lists returned by functions and methods
// in the file and indexer packages are sorted this way unless undefined.
func Intersect(a, b []*Page) (c []*Page) {
	for aI, bI := 0, 0; aI < len(a) && bI < len(b); {
		aAddy := a[aI].Address()
		bAddy := b[bI].Address()
		switch {
		case aAddy < bAddy:
			aI++
		case aAddy > bAddy:
			bI++
		default:
			c = append(c, a[aI])
			aI++
			bI++
		}
	}
	return c
}

// Subtract returns a list of pages in a but not in b. Its output is undefined
// unless both are sorted by address. Lists returned by functions and methods
// in the file and indexer packages are sorted this way unless undefined.
func Subtract(a, b []*Page) (c []*Page) {
	aI := 0
	for bI := 0; aI < len(a) && bI < len(b); {
		aAddy := a[aI].Address()
		bAddy := b[bI].Address()
		switch {
		case aAddy > bAddy:
			// don't know what to do with a yet, bring b forward
			bI++
		case aAddy < bAddy:
			// keep a, bring a forward
			c = append(c, a[aI])
			aI++
		default:
			// discard a, bring both forward
			aI++
			bI++
		}
	}
	// keep any remaining a's, have run out of b's
	c = append(c, a[aI:]...)
	return c
}

// Add returns a list of pages in either a or b. Its output is undefined
// unless both are sorted by address. Lists returned by functions and methods
// in the file and indexer packages are sorted this way unless undefined.
func Add(a, b []*Page) []*Page {
	c := make([]*Page, 0, len(a)+len(b))
	for len(a) > 0 && len(b) > 0 {
		aAddy := a[0].Address()
		bAddy := b[0].Address()
		switch {
		case aAddy < bAddy:
			c = append(c, a[0])
			a = a[1:]
		case aAddy > bAddy:
			c = append(c, b[0])
			b = b[1:]
		default:
			c = append(c, a[0])
			a = a[1:]
			b = b[1:]
		}
	}
	if len(a) > 0 {
		c = append(c, a...)
	} else if len(b) > 0 {
		c = append(c, b...)
	}
	// don't bother to copy c to a correctly sized slice. many duplicates are
	// unlikely in practice.
	return c
}
