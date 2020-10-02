package build_test

import (
	"testing"

	"github.com/go-test/deep"
	"github.com/shogg/edifact/internal/build"
)

func TestSegmentSelector(t *testing.T) {

	tests := []struct {
		structTag string
		expected  build.SegmentSelector
	}{
		{
			`DTM+4|5`,
			build.SegmentSelector{
				Tag:    "DTM",
				Params: []build.SegmentSelectorParam{{0, 0, "4|5"}},
			},
		},
		{
			`SG1/RFF+100`,
			build.SegmentSelector{
				Path:   "SG1/",
				Tag:    "RFF",
				Params: []build.SegmentSelectorParam{{0, 0, "100"}},
			},
		},
		{
			`SG1/RFF+300|301`,
			build.SegmentSelector{
				Path:   "SG1/",
				Tag:    "RFF",
				Params: []build.SegmentSelectorParam{{0, 0, "300|301"}},
			},
		},
		{
			`SG10/SG17/LIN+?`,
			build.SegmentSelector{
				Path:  "SG10/SG17/",
				Tag:   "LIN",
				Value: build.ValueComponent{0, 0},
			},
		},
		{
			`SG10/SG17/LIN+++?:`,
			build.SegmentSelector{
				Path:  "SG10/SG17/",
				Tag:   "LIN",
				Value: build.ValueComponent{2, 0},
			},
		},
		{
			`SG10/SG17/QTY++:?`,
			build.SegmentSelector{
				Path:  "SG10/SG17/",
				Tag:   "QTY",
				Value: build.ValueComponent{1, 1},
			},
		},
	}

	for _, test := range tests {
		sb := build.ParseSegmentSelector(test.structTag)
		if diff := deep.Equal(test.expected, sb); diff != nil {
			t.Error(diff)
		}
	}
}
