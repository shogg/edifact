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
DTM+4:20060620:102'
RFF+100+1'
RFF+301+2'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++product A:SA'
QTY+2:1000'
LIN+2++product B:SA'
QTY+3:1000'
CNT+2:1'
UNT+9+1'

UNH+1+DESADV:D:96A:UN'
BGM+220+B10001'
DTM+4:20060620:102'
RFF+100+1'
RFF+301+2'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++product A:SA'
QTY+2:1000'
LIN+2++product B:SA'
QTY+3:1000'
CNT+2:1'
UNT+9+1'
`

func TestUnmarshal(t *testing.T) {

	type Item struct {
		ItemNr      int    `edifact:"SG10/SG17/LIN+?"`
		Description string `edifact:"SG10/SG17/LIN+++?"`
		Quantity    int    `edifact:"SG10/SG17/QTY+?"`
	}
	type Message struct {
		Date       time.Time `edifact:"DTM+4|5"`
		DeliveryNr string    `edifact:"SG1/RFF+100+?"`
		OrderNr    int       `edifact:"SG1/RFF+300|301+?"`
		Items      []Item    `edifact:"SG10/SG17"`
	}

	document := strings.NewReader(ediMessage)
	var ediData []*Message
	if err := edifact.Unmarshal(document, &ediData); err != nil {
		t.Fatal(err)
	}

	if len(ediData) == 0 {
		t.Error("data expected")
	}
	if ediData[0].Items[0].ItemNr != 1 {
		t.Error("LIN+1 expected")
	}
	if ediData[0].Items[0].Description != "product A" {
		t.Error("product A expected")
	}
	if ediData[0].Items[1].ItemNr != 2 {
		t.Error("quantity 1 expected")
	}
	if ediData[0].Items[1].Description != "product B" {
		t.Error("product B expected")
	}
}
