package build

import (
	"strings"
)

// SegmentSelector data structure representing struct tag `edifact:""`.
type SegmentSelector struct {
	Path   string
	Tag    string
	Params []SegmentSelectorParam
	Value  ValueComponent
}

// SegmentSelectorParam selection parameters.
type SegmentSelectorParam struct {
	Elem  int
	Comp  int
	Value string
}

// ValueComponent segment component holding the value.
type ValueComponent struct {
	Elem int
	Comp int
}

// ParseSegmentSelector parses a struct tag.
func ParseSegmentSelector(structTag string) SegmentSelector {
	splitted := strings.SplitN(structTag, "+", 2)

	pathEnd := strings.LastIndexByte(splitted[0], '/')
	path := ""
	tag := splitted[0]
	if pathEnd >= 0 {
		path = splitted[0][:pathEnd+1]
		tag = splitted[0][pathEnd+1 : pathEnd+1+3]
	} else {
		tag = splitted[0]
	}

	var params []SegmentSelectorParam
	var value ValueComponent

	if len(splitted) == 2 {
		params, value = parseParamsAndValue(splitted[1])
	}

	return SegmentSelector{
		Path:   path,
		Tag:    tag,
		Params: params,
		Value:  value,
	}
}

func parseParamsAndValue(s string) ([]SegmentSelectorParam, ValueComponent) {

	var params []SegmentSelectorParam
	var value ValueComponent

	elems := strings.Split(s, "+")
	for i, e := range elems {

		comps := strings.Split(e, ":")
		for j, c := range comps {
			if len(c) == 0 {
				continue
			}
			if c == "?" {
				value = ValueComponent{
					Elem: i,
					Comp: j,
				}
				continue
			}

			params = append(params, SegmentSelectorParam{
				Elem:  i,
				Comp:  j,
				Value: c,
			})
		}
	}

	return params, value
}
