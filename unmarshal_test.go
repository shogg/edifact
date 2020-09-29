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
LIN+1++Produkt Schrauben:SA'
QTY+1:1000'
UNS+S'
CNT+2:1'
UNT+9+1'
UNZ+1+1234567'
`

func TestUnmarshal(t *testing.T) {

	type TestMessage struct {
		Datum      time.Time `edifact:"DTM+4|5"`
		AuftragNr  string    `edifact:"SG1/RFF+100"`
		BestellNr  string    `edifact:"SG1/RFF+300|301"`
		Positionen []struct {
			PositionNr  int `edifact:"SG10/LIN+?"`
			Bezeichnung int `edifact:"SG10/LIN+++?:"`
			Anzahl      int `edifact:"SG10/QTY+?:"`
		}
	}

	document := strings.NewReader(ediMessage)
	ediData := TestMessage{}

	if err := edifact.Unmarshal(document, &ediData); err != nil {
		t.Fatal(err)
	}

	if len(ediData.Positionen) == 0 {
		t.Error("expected positions")
	}
}
