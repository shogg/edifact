package spec_test

import (
	"testing"

	"github.com/shogg/edifact/spec"
)

func TestSegment(t *testing.T) {

	seg := spec.Segment("UNH+1+?+ORDERS:D:96A:UN'")

	tests := []struct {
		e, c     int
		expected string
	}{
		{0, 0, "UNH"},
		{0, 1, ""},
		{1, 0, "1"},
		{1, 1, ""},
		{2, 0, "?+ORDERS"},
		{2, 1, "D"},
		{2, 3, "UN"},
	}

	for _, test := range tests {
		actual := seg.Elem(test.e).Comp(test.c)
		if actual != test.expected {
			t.Errorf(`"%s" expected, was "%s"`, test.expected, actual)
		}
	}
}

func BenchmarkSegment(b *testing.B) {

	seg := spec.Segment("UNH+1+?+ORDERS:D:96A:UN'")
	for i := 0; i < b.N; i++ {
		_ = seg.Elem(2).Comp(3)
	}
}
