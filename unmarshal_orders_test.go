package edifact_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/shogg/edifact"
)

var ordersMessage = `
UNA:+.? '
UNB+UNOC:3+3020400000100:14+5909000580916:14+180607:1532+164702918+ORDERS'
UNG+ORDERS'
UNH+1+ORDERS:D:96A:UN:EAN008'
BGM+220::9+0100623447+9'
DTM+137:20180607:102'
DTM+2:20180605:102'
FTX+PUR+3++ For the appointment, please contact us on 04 79 68 55 40 (Fax 04:79 68 55 41)’
FTX+PUR+3++ STORE CUSTOMER ORDER:CP01’
RFF+UC:0000979343'
NAD+BY+3020400000100::9'
NAD+SU+5909000580916::9'
NAD+DP+3020400235403::9'
NAD+UC+++COIGNIERES+ROUTE NATIONALE 10:RUE DES LOUVERIES+COIGNIERES++78310+FR'
CTA+GR'
COM+01 30 49 24 60:TE'
CUX+2:EUR:9+3:EUR:4'
LIN+00001++3663602943105:EN'
PIA+5+3663602943105:SA'
IMD+F++:::POT VERT 7X7X200 OMBRONE CL4'
QTY+21:80:EA'
QTY+52:1:EA'
FTX+PUR+++POT VERT 7X7X200 OMBRONE CL4--3663602943105'
PRI+AAA:3.31::NTP:1:EA'
UNS+S'
CNT+2:1'
UNT+24+1'
UNE+ORDERS'
UNZ+1+164702918'`

type Order struct {
	OrderNo      string       `edifact:"BGM+220+?"`
	OrderDate    time.Time    `edifact:"DTM+2:?"`
	DeliveryDate time.Time    `edifact:"DTM+137:?"`
	CustomerNo   string       `edifact:"SG2/NAD+BY+?"`
	Items        []*OrderItem `edifact:"SG25"`
}

type OrderItem struct {
	ItemNo   int  `edifact:"SG25/LIN+?"`
	Quantity int  `edifact:"SG25/QTY+21:?"`
	GTIN     GTIN `edifact:"SG25/LIN+++*:EN"`
}

type GTIN string

func (gtin *GTIN) UnmarshalEdifact(data []byte) error {
	*gtin = GTIN(data)
	return nil
}

func TestUnmarshalOrders(t *testing.T) {

	document := strings.NewReader(ordersMessage)
	var orders []*Order
	if err := edifact.Unmarshal(document, &orders); err != nil {
		t.Error(err)
	}

	data, err := json.MarshalIndent(orders, "", "\t")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(data))
}
