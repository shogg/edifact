package build

import (
	"testing"

	"github.com/go-test/deep"
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
				params: []segmentSelectorParam{{0, 0, "4|5"}},
			},
		},
		{
			`SG1/RFF+100`,
			segmentSelector{
				path:   "SG1/",
				tag:    "RFF",
				params: []segmentSelectorParam{{0, 0, "100"}},
			},
		},
		{
			`SG1/RFF+300|301`,
			segmentSelector{
				path:   "SG1/",
				tag:    "RFF",
				params: []segmentSelectorParam{{0, 0, "300|301"}},
			},
		},
		{
			`SG10/SG17/LIN+?`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{0, 0},
			},
		},
		{
			`SG10/SG17/LIN+++?:`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{2, 0},
			},
		},
		{
			`SG10/SG17/LIN+++*`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{2, -1},
			},
		},
		{
			`SG10/SG17/LIN+++*:EN`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{2, -1},
			},
		},
		{
			`SG10/SG17/LIN+++:*`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "LIN",
				value: valueComponent{2, -1},
			},
		},
		{
			`SG10/SG17/QTY++:?`,
			segmentSelector{
				path:  "SG10/SG17/",
				tag:   "QTY",
				value: valueComponent{1, 1},
			},
		},
	}

	for _, test := range tests {
		sb := parseSegmentSelector(test.structTag)
		if diff := deep.Equal(test.expected, sb); diff != nil {
			t.Error(diff)
		}
	}
}
