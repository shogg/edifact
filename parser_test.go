package edifact

import (
	"fmt"
	"strings"
	"testing"

	"github.com/shogg/edifact/spec"
)

var ediMessage1 = `
UNA:+.? '
UNB+UNOC:3+Senderkennung+Empfaengerkennung+060620:0931+1++1234567'
UNH+1+ORDERS:D:96A:UN'
BGM+220+B10001'
DTM+4:20060620:102'
NAD+BY+++Bestellername+Strasse+Stadt++23436+xx'
LIN+1++Produkt Schrauben:SA'
QTY+1:1000'
UNS+S'
CNT+2:1'
UNT+9+1'
UNZ+1+1234567'
`

func TestParser(t *testing.T) {

	p := newParser(strings.NewReader(ediMessage1))
	h := &testHandler{}
	if err := p.parse(h); err != nil {
		t.Fatal(err)
	}
}

type testHandler struct{}

func (h *testHandler) Handle(sgs []*spec.Node, s Segment) error {
	fmt.Println(formatSegmentGroups(sgs), s)
	return nil
}

func formatSegmentGroups(sgs []*spec.Node) string {
	var buf strings.Builder
	for i, n := range sgs {
		if i != 0 {
			buf.WriteByte('/')
		}
		buf.WriteString(n.Tag)
	}
	return buf.String()
}
