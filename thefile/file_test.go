package thefile

import (
	"reflect"
	"testing"
)

func TestAddress(t *testing.T) {
	buf := []byte(`----one dot one
----one dot two

----two dot one

two content
lines

`)
	pages, _ := pagesFrom(buf)
	want := []int{1, 4}
	var got []int
	for _, page := range pages {
		got = append(got, page.Address())
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("\nwant: %v\ngot:  %v\n", want, got)
	}
}

var ps []*Page

func BenchmarkPages(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var err error
		ps, err = Pages()
		if err != nil {
			b.Error(err)
			return
		}
	}
}
