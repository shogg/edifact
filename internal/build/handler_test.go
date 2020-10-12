package build

import (
	"testing"

	"github.com/shogg/edifact/spec"
)

func TestHandleStruct(t *testing.T) {

	model := spec.Msg("Test", spec.S("RFF", spec.C, 1))

	value := struct {
		RFF int `edifact:"RFF++:?"`
	}{}

	h := &Handler{Target: &value}
	seg := spec.Segment(`RFF+bla+:1:fasel'`)
	if err := h.Handle(model.FirstChild, seg); err != nil {
		t.Error(err)
	}

	if value.RFF != 1 {
		t.Error("RFF 1 expected")
	}
}
