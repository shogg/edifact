package build

import (
	"testing"

	"github.com/shogg/edifact/spec"
)

func TestUnmarshalSimple(t *testing.T) {

	type simple struct {
		String   string
		Bool     bool
		Int      int
		Int64    int64
		Uintptr  uintptr
		Float    float64
		Complex  complex128
		TypedPtr *int
	}

	seg := spec.Segment(`RFF++:1'`)
	comp := ValueComponent{2, 1}

	var s simple
	if err := unmarshal(&s, seg, comp); err != nil {
		t.Fatal(err)
	}

	if s.String != "1" {
		t.Error("simple.Int: \"1\" expected, was", s.String)
	}
	if s.Bool != true {
		t.Error("simple.Bool: true expected, was", s.Bool)
	}
	if s.Int != 1 {
		t.Error("simple.Int: 1 expected, was", s.Int)
	}
	if s.Int64 != 1 {
		t.Error("simple.Int64: 1 expected, was", s.Int64)
	}
	if s.Uintptr != 1 {
		t.Error("simple.Uintptr: 1 expected, was", s.Uintptr)
	}
	if s.Float != 1.0 {
		t.Error("simple.Float: 1.0 expected, was", s.Uintptr)
	}
	if s.Complex != 1+0i {
		t.Error("simple.Complex: 1+0i expected, was", s.Complex)
	}
	if s.TypedPtr == nil {
		t.Error("simple.TypedPtr: 1 expected, was", *s.TypedPtr)
	}
}
