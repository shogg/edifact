package edifact_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shogg/edifact"
)

var ediMessage = `
UNA:+.? '
UNB+UNOC:3+sender+receiver+060620:0931+1++1234567'

UNH+1+DESADV:D:96A:UN'
BGM+220+B10001'
DTM+17:20060620:102'
RFF+ON+1'
RFF+DQ+2'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++Beck?'s:SA'
QTY+12:10'
LIN+2++Chantr??'
QTY+12:20'
CNT+2:1'
UNT+9+1'

UNH+1+DESADV:D:96A:UN'
BGM+220+B10001'
DTM+17:20060620:102'
RFF+ON+1'
RFF+DQ+2'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++Beck?'s:SA'
QTY+12:10'
LIN+2++Chantr??'
QTY+12:20'
CNT+2:1'
UNT+9+1'
`

type Message struct {
	Date       time.Time `edifact:"DTM+17|18"`
	DeliveryNr string    `edifact:"SG1/RFF+VN|DQ+?"`
	OrderNr    int       `edifact:"SG1/RFF+ON+?"`
	Items      []Item    `edifact:"SG10/SG17"`
}
type Item struct {
	Message     *Message
	ItemNr      int    `edifact:"SG10/SG17/LIN+?"`
	Description string `edifact:"SG10/SG17/LIN+++?"`
	Quantity    int    `edifact:"SG10/SG17/QTY+12:?"`
}

func TestUnmarshal(t *testing.T) {

	ediData, err := unmarshal()
	if err != nil {
		t.Fatal(err)
	}

	if len(ediData) == 0 {
		t.Error("data expected")
	}
	if ediData[0].Items[0].ItemNr != 1 {
		t.Error("LIN+1 expected")
	}
	desc0 := ediData[0].Items[0].Description
	if desc0 != "Beck's" {
		t.Error("Beck's expected, was", desc0)
	}
	if ediData[0].Items[1].ItemNr != 2 {
		t.Error("quantity 1 expected")
	}
	desc1 := ediData[0].Items[1].Description
	if desc1 != "Chantr?" {
		t.Error("Chantr? expected, was", desc1)
	}
}

func BenchmarkUnmarshal(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = unmarshal()
	}
}

func unmarshal() ([]*Message, error) {

	document := strings.NewReader(ediMessage)
	var ediData []*Message
	if err := edifact.Unmarshal(document, &ediData); err != nil {
		return nil, err
	}

	return ediData, nil
}
