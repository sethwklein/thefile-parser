package thefile

import (
	"testing"
)

var list []*Page

func BenchmarkIntersect(b *testing.B) {
	pages, err := Pages()
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 100; j++ {
			list = Intersect(pages[j:j*10], pages[j+100:(j*10)+100])
		}
	}
}

// BUG(sk): these are pretty half hearted tests

func TestSubtract(t *testing.T) {
	pages, err := Pages()
	if err != nil {
		t.Error(err)
		return
	}
	super := pages[:20]
	sub := pages[5:15]
	outer := Subtract(super, sub)
	for _, p := range outer {
		for _, q := range sub {
			if p.Address() == q.Address() {
				t.Errorf("found subtracted item\nsuper: %v\nsub:  %v\nouter: %v\n", super, sub, outer)
				return
			}
		}
	}
}

func TestIntersect(t *testing.T) {
	pages, err := Pages()
	if err != nil {
		t.Error(err)
		return
	}
	first := pages[:10]
	second := pages[5:15]
	both := Intersect(first, second)
	for i, p := range both {
		if p.Address() != pages[i+5].Address() {
			t.Error("wrong thing here")
		}
	}
}
