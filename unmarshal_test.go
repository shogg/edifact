package edifact_test

import (
	"strings"
	"testing"
	"time"

	"github.com/shogg/edifact"
)

var ediMessage = `
UNA:+.? '
UNB+UNOC:3+Senderkennung+Empfaengerkennung+060620:0931+1++1234567'
UNH+1+ORDERS:D:96A:UN'
BGM+220+B10001'
DTM+4:20060620:102'
NAD+BY+++Bestellername+Strasse+Stadt++23436+xx'
CPS+'
LIN+1++product bolts:SA'
QTY+1:1000'
LIN+2++product nuts:SA'
QTY+1:1000'
CNT+2:1'
UNT+9+1'
`

func TestUnmarshal(t *testing.T) {

	type TestMessage struct {
		Datum      time.Time `edifact:"DTM+4|5"`
		AuftragNr  string    `edifact:"SG1/RFF+100"`
		BestellNr  string    `edifact:"SG1/RFF+300|301"`
		Positionen []struct {
			PositionNr  int    `edifact:"SG10/SG17/LIN+?"`
			Bezeichnung string `edifact:"SG10/SG17/LIN+++?:"`
			Anzahl      int    `edifact:"SG10/SG17/QTY+?:"`
		} `edifact:"SG10/SG17"`
	}

	document := strings.NewReader(ediMessage)
	var ediData []*TestMessage
	if err := edifact.Unmarshal(document, &ediData); err != nil {
		t.Fatal(err)
	}

	if len(ediData) == 0 {
		t.Error("data expected")
	}
	if ediData[0].Positionen[0].PositionNr != 1 {
		t.Error("LIN+1 expected")
	}
	if ediData[0].Positionen[0].Bezeichnung != "product bolts" {
		t.Error("product bolts expected")
	}
	if ediData[0].Positionen[1].PositionNr != 2 {
		t.Error("LIN+1 expected")
	}
	if ediData[0].Positionen[1].Bezeichnung != "product nuts" {
		t.Error("product bolts expected")
	}
}
