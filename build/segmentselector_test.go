package build

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSegmentSelector(t *testing.T) {

	tests := []struct {
		structTag string
		expected  segmentSelector
	}{
		{
			`DTM+4|5`,
			segmentSelector{
				tag:    "DTM",
				params: []segmentSelectorParam{{1, 0, "4|5"}},
			},
		},
		{
			`SG1/RFF+100`,
			segmentSelector{
				path:   "SG1/",
				tag:    "RFF",
				params: []segmentSelectorParam{{1, 0, "100"}},
			},
		},
		{
			`SG1/RFF+300|301`,
			segmentSelector{
				path:   "SG1/",
				tag:    "RFF",
				params: []segmentSelectorParam{{1, 0, "300|301"}},
			},
		},
		{
			`SG10/SG17/LIN+?`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{1, 0},
			},
		},
		{
			`SG10/SG17/LIN+++?:`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{3, 0},
			},
		},
		{
			`SG10/SG17/LIN+++*`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{3, -1},
			},
		},
		{
			`SG10/SG17/LIN+++*:EN`,
			segmentSelector{
				path:   "SG10/SG17/",
				tag:    "LIN",
				params: []segmentSelectorParam{{3, 1, "EN"}},
				value:  valueComponent{3, -1},
			},
		},
		{
			`SG10/SG17/LIN+++:*`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{3, -1},
			},
		},
		{
			`SG10/SG17/QTY++:?`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "QTY",
				value: valueComponent{2, 1},
			},
		},
	}

	for _, test := range tests {
		sb := parseSegmentSelector(test.structTag)
		assert.Equal(t, test.expected, sb)
	}
}
