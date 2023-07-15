package parse_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shogg/edifact/parse"
	"github.com/shogg/edifact/spec"
)

var ediMessage1 = `
UNA:+.? '
UNB+UNOC:3+Senderkennung+Empfaengerkennung+060620:0931+1++1234567'
UNH+1+ORDERS:D:96A:UN'
BGM+220+B10001'
DTM+4:20060620:102'
RFF+1+1++'
DTM+4:20060620:102'
RFF+2+2++'
NAD+BY+++Bestellername+Strasse+Stadt++23436+xx'
LIN+1++Produkt Schrauben:SA'
QTY+1:1000'
CNT+2:1'
UNT+9+1'
`

func TestMissingSegmentTerminator(t *testing.T) {

	tests := []string{
		"DTA",
		"DTA'DTA",
		"DTA'\nDTA",
		"DTA'\n\nDTA",
	}

	h := testHandler{}
	for _, test := range tests {
		err := parse.Parse(strings.NewReader(test), &h)
		if err == nil {
			t.Error("missing segment terminator expected")
		} else {
			fmt.Println(err)
		}
	}
}

func TestReleaseChar(t *testing.T) {

	tests := []string{
		"DTA'\n",
		"DTA\n'",
		"DTA?''",
		"DTA??'",
		"DTA???''",
		"DTA?'?''",
		"DTA?'??'",
	}

	h := testHandler{}
	for _, test := range tests {
		if err := parse.Parse(strings.NewReader(test), &h); err != nil {
			t.Error(err)
		}
	}
}

func TestParser(t *testing.T) {
	h := testHandler{}
	if err := parse.Parse(strings.NewReader(ediMessage1), &h); err != nil {
		t.Fatal(err)
	}
}

type testHandler struct{}

func (h *testHandler) Handle(node *spec.Node, s spec.Segment) error {
	fmt.Printf("%s%s\n", formatSegmentGroups(node), s)
	return nil
}

func formatSegmentGroups(node *spec.Node) string {
	var buf strings.Builder
	for _, sg := range node.SegmentGroups() {
		buf.WriteString(sg.Tag)
		buf.WriteByte('/')
	}
	return buf.String()
}
