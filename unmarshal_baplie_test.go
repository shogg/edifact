package edifact_test

import (
	"strings"
	"testing"

	"github.com/shogg/edifact"
	"github.com/shogg/edifact/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBaplie(t *testing.T) {
	spec.Add("BAPLIE", BAPLIE)
	document := strings.NewReader(testMsgWith1Container)
	var vd VesselVisitDetails
	err := edifact.Unmarshal(document, &vd)
	require.NoError(t, err)
	assert.Equal(t, 3, len(vd.Containers[0].DangerousGoods))
	assert.Equal(t, "9", vd.Containers[0].DangerousGoods[0].IMOClass)
	assert.Equal(t, "8", vd.Containers[0].DangerousGoods[1].IMOClass)
	assert.Equal(t, "3.1", vd.Containers[0].DangerousGoods[2].IMOClass)
}

var testMsgWith1Container = `
	UNB+UNOA:2+SEACOS+PUBLIC+210118:1038+31033627+++++UNKNOWN'
	UNH+SEACOSA3627+BAPLIE:D:95B:UN:SMDG22'
	BGM++R0C55E2211C418+9'
	DTM+137:2101181038:201'
	TDT+20+010E+++:172:20+++A8VP5:103:ZZZ:'
	LOC+5+GBLGP:139:6'
	LOC+61+DEHAM:139:6'
	DTM+133+test'
	LOC+147+0371582::5'
	FTX+HAN+++OND'
	FTX+AAY+++202205091336UTC'
	MEA+WT++KGM:26700'
	MEA+VGM++KGM:26700'
	TMP+2+-18:CEL'
	LOC+9+DEHAM:139:6+CTT:TER:306'
	LOC+11+SGSIN:139:6'
	RFF+BN:2697489840'
	EQD+CN+CSNU1487660+22G1:102:5++2+5'
	NAD+CA+OOL:172:ZZZ'
	DGS+IMD+9+3082++3+F-AS-F'
	DGS+IMD+8+3133++4+F-AS-F'
	FTX+AAD+++PHENOXYBENZYL-(1R -TRANS) (-2,2-DIMETHYL - 3(2-METHYLPROP-1-ENYL) CYCL:220.0'
	DGS+IMD+3.1+1263+:+++++:+'
	FTX+AAD+++PHENOXYBENZYL-(1R -TRANS) (-2,2-DIMETHYL - 3(2-METHYLPROP-1-ENYL) CYCL:220.0'
	UNT+18142+SEACOSA3627'
	UNZ+1+31033627'`

type VesselVisitDetails struct {
	Containers []Container `edifact:"SG2"`
}

type Container struct {
	ReferenceNumber    string          `edifact:"SG2/RFF+BN:?"`
	IsoCode            string          `edifact:"SG2/SG3/EQD+CN+?+?"`
	ContainerNumber    string          `edifact:"SG2/SG3/EQD+CN+?"`
	PortOfLoading      string          `edifact:"SG2/LOC+9+?"`
	PortOfDischarge    string          `edifact:"SG2/LOC+11+?"`
	StowageLocation    string          `edifact:"SG2/LOC+147+?"`
	Operator           string          `edifact:"SG2/SG3/NAD+CA+?"`
	FullEmptyIndicator string          `edifact:"SG2/SG3/EQD+CN+?+?+++?"`
	GrossWeight        string          `edifact:"SG2/MEA+VGM++KGM:?"`
	Temperature        string          `edifact:"SG2/TMP+2+?"`
	DangerousGoods     []DangerousGood `edifact:"SG2/SG4"`
}

type DangerousGood struct {
	IMOClass string `edifact:"SG2/SG4/DGS+IMD+?"`
}

var BAPLIE = spec.Msg("BAPLIE",
	spec.S("UNB", spec.C, 1),
	spec.S("UNH", spec.M, 1),
	spec.S("BGM", spec.M, 1),
	spec.S("DTM", spec.M, 1),
	spec.S("RFF", spec.C, 1),
	spec.S("NAD", spec.C, 9),
	spec.SG("SG1", spec.C, 1,
		spec.S("TDT", spec.M, 1),
		spec.S("LOC", spec.M, 9),
		spec.S("DTM", spec.M, 99),
		spec.S("RFF", spec.C, 1),
		spec.S("FTX", spec.C, 1),
	),
	spec.SG("SG2", spec.C, 9999,
		spec.S("LOC", spec.M, 1),
		spec.S("GID", spec.C, 1),
		spec.S("GDS", spec.C, 9),
		spec.S("FTX", spec.C, 9),
		spec.S("MEA", spec.M, 9),
		spec.S("DIM", spec.C, 9),
		spec.S("TMP", spec.C, 1),
		spec.S("RNG", spec.C, 1),
		spec.S("LOC", spec.C, 9),
		spec.S("RFF", spec.M, 9),
		spec.SG("SG3", spec.C, 9,
			spec.S("EQD", spec.M, 1),
			spec.S("EQA", spec.C, 9),
			spec.S("NAD", spec.C, 1),
		),
		spec.SG("SG4", spec.C, 999,
			spec.S("DGS", spec.M, 1),
			spec.S("FTX", spec.C, 1),
		),
	),
	spec.S("UNT", spec.M, 1),
	spec.S("UNZ", spec.M, 1),
)
