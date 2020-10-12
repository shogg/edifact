package build

import (
	"strings"

	"github.com/shogg/edifact/spec"
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
					Elem: i + 1,
					Comp: j,
				}
				continue
			}

			params = append(params, SegmentSelectorParam{
				Elem:  i + 1,
				Comp:  j,
				Value: c,
			})
		}
	}

	return params, value
}

// Select returns a seqment component value.
func (sel SegmentSelector) Select(seg spec.Segment) string {
	if sel.Value.Elem == 0 {
		return string(seg)
	}
	return seg.Elem(sel.Value.Elem).Comp(sel.Value.Comp)
}

// Matches returns true if segment group path, segment tag and segment selector parameters matches.
func (sel SegmentSelector) Matches(node *spec.Node, seg spec.Segment) bool {

	if sel.Tag != seg.Tag() {
		return false
	}
	if sel.Path != node.Path() {
		return false
	}

	matches := true
	for _, param := range sel.Params {
		matches = matches && param.matches(seg)
	}
	return matches
}

// matches returns true if a segment component value matches the param.
func (param SegmentSelectorParam) matches(seg spec.Segment) bool {
	comp := seg.Elem(param.Elem).Comp(param.Comp)
	candidates := strings.Split(param.Value, "|")
	for _, c := range candidates {
		if comp == c {
			return true
		}
	}
	return false
}
