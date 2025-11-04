package spec_test

import (
	"testing"

	"github.com/shogg/edifact/spec"
)

func TestReleaseChar(t *testing.T) {

	tests := []string{
		"UNH+1+ORDERS:D:96A:UN'",
		"UNH?++1+ORDERS:D:96A:UN'",
		"UNH??+1+ORDERS:D:96A:UN'",
		"UNH???++1+ORDERS:D:96A:UN'",
		"UNH?+?++1+ORDERS:D:96A:UN'",
		"UNH+1+ORDERS?::D:96A:UN'",
		"UNH+1+ORDERS??:D:96A:UN'",
		"UNH+1+ORDERS???::D:96A:UN'",
		"UNH+1+ORDERS?:?::D:96A:UN'",
		"UNH+1+ORDER?'S:D:96A:UN'",
		"UNH+1+ORDER???'S:D:96A:UN'",
	}

	for i, test := range tests {
		got := spec.Segment(test).Comp(2).Elem(3)
		if got != "UN" {
			t.Errorf("%d: wanted \"UN\" got \"%s\"", i, got)
		}
	}
}

func TestSegment(t *testing.T) {

	seg := spec.Segment("UNH+1+?+?+ORDER?'S???:?::D:96A:UN'")

	tests := []struct {
		c, e     int
		expected string
	}{
		{0, 0, "UNH"},
		{0, 1, ""},
		{1, 0, "1"},
		{1, 1, ""},
		{2, 0, "++ORDER'S?::"},
		{2, 1, "D"},
		{2, 3, "UN"},
		{3, 0, ""},
	}

	for i, test := range tests {
		actual := seg.Comp(test.c).Elem(test.e)
		if actual != test.expected {
			t.Errorf(`"%d: %s" expected, was "%s"`, i, test.expected, actual)
		}
	}
}

func BenchmarkSegment(b *testing.B) {

	seg := spec.Segment("UNH+1+?+ORDERS:D:96A:UN'")
	for i := 0; i < b.N; i++ {
		_ = seg.Comp(2).Elem(3)
	}
}

func BenchmarkTag(b *testing.B) {

	seg := spec.Segment("UNH+1+?+ORDERS:D:96A:UN'")
	for i := 0; i < b.N; i++ {
		_ = seg.Tag()
	}
}
