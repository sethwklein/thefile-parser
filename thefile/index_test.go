package thefile

import (
	"testing"

	"sethwklein.net/thefile/thefile"
)

var index *Index

func BenchmarkNew(b *testing.B) {
	pages, err := thefile.Pages()
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		index = NewIndex(pages)
	}
}
