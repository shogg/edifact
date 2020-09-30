package spec_test

import (
	"testing"

	"github.com/shogg/edifact/spec"
)

func TestTransition(t *testing.T) {

	unh, err := spec.DESADV.Transition("UNH")
	if err != nil {
		t.Fatal(err)
	}
	bgm, err := unh.Transition("BGM")
	if err != nil {
		t.Fatal(err)
	}
	rff, err := bgm.Transition("RFF")
	if err != nil {
		t.Fatal(err)
	}
	if rff == nil || rff.Tag != "RFF" {
		t.Error("RFF expected")
	}
	nad, err := rff.Transition("NAD")
	if err != nil {
		t.Fatal(err)
	}
	if nad == nil || nad.Tag != "NAD" {
		t.Error("NAD expected")
	}
	_, err = nad.Transition("DTM")
	if err == nil {
		t.Error("error expected")
	}
}
