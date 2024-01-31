package edifact_test

import (
	"strings"
	"testing"

	"github.com/shogg/edifact"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var ediMessageRFF = `
UNA:+.? '
UNB+UNOC:3+sender+receiver+060620:0931+1++1234567'
UNG+DESADV'

UNH+1+DESADV:D:96A:UN'
BGM+220+B10001'
DTM+17:20060620:102'
RFF+ON+1'
RFF+DQ+2'
RFF+ON+3'
RFF+DQ+4'
NAD+BY+++name+street+city++23436+xx'
CPS+'
LIN+1++Beck?'s:SA'
QTY+12:10'
LIN+2++Chantr??'
QTY+12:20'
CNT+2:1'
UNT+9+1'

UNE+'
UNZ+'
`

type MessageRFF struct {
	RFFs   []string `edifact:"SG1/RFF"`
	RFFsON []int    `edifact:"SG1/RFF+ON+?"`
	RFFsDQ []int    `edifact:"SG1/RFF+DQ+?"`
}

func TestUnmarshalRFFs(t *testing.T) {

	document := strings.NewReader(ediMessageRFF)
	var messages []*MessageRFF
	err := edifact.Unmarshal(document, &messages)
	require.NoError(t, err)

	assert.Len(t, messages, 1)
	assert.Len(t, messages[0].RFFs, 4)
	assert.Len(t, messages[0].RFFsON, 2)
	assert.Len(t, messages[0].RFFsDQ, 2)
	assert.Equal(t, []string{"RFF+ON+1'", "RFF+DQ+2'", "RFF+ON+3'", "RFF+DQ+4'"}, messages[0].RFFs)
	assert.Equal(t, []int{1, 3}, messages[0].RFFsON)
	assert.Equal(t, []int{2, 4}, messages[0].RFFsDQ)
}
