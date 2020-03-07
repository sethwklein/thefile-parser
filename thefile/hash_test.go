package thefile

import (
	"testing"
)

var hash []byte

func BenchmarkHash(b *testing.B) {
	pages, err := Pages()
	if err != nil {
		b.Error(err)
		return
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1; j++ {
			for _, page := range pages {
				hash = page.Hash()
			}
		}
	}
}
