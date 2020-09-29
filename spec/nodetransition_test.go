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
	sg1, err := bgm.Transition("RFF")
	if err != nil {
		t.Fatal(err)
	}
	if sg1 == nil || sg1.Tag != "SG1" {
		t.Error("SG1 expected")
	}
	rff, err := sg1.Transition("RFF")
	if err != nil {
		t.Fatal(err)
	}
	if rff == nil || rff.Tag != "RFF" {
		t.Error("RFF expected")
	}
	sg2, err := rff.Transition("NAD")
	if err != nil {
		t.Fatal(err)
	}
	if sg2 == nil || sg2.Tag != "SG2" {
		t.Error("SG2 expected")
	}
	nad, err := sg2.Transition("NAD")
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
