package build

import (
	"strings"

	"github.com/shogg/edifact/spec"
)

// segmentSelector represents struct tag `edifact:""`.
type segmentSelector struct {
	path   string
	tag    string
	params []segmentSelectorParam
	value  valueComponent
}

// segmentSelectorParam segment filter parameter.
type segmentSelectorParam struct {
	comp  int
	elem  int
	value string
}

// valueComponent value position in the segment.
type valueComponent struct {
	comp int
	elem int
}

// parseSegmentSelector parses a struct tag.
func parseSegmentSelector(structTag string) segmentSelector {
	splitted := strings.SplitN(structTag, "+", 2)

	pathEnd := strings.LastIndexByte(splitted[0], '/')
	path := ""
	tag := splitted[0]
	if pathEnd >= 0 {
		path = splitted[0][:pathEnd+1]
		tag = splitted[0][pathEnd+1:]
	}

	var params []segmentSelectorParam
	var value valueComponent

	if len(splitted) == 2 {
		params, value = parseParamsAndValue(splitted[1])
	}

	return segmentSelector{
		path:   path,
		tag:    tag,
		params: params,
		value:  value,
	}
}

func parseParamsAndValue(s string) ([]segmentSelectorParam, valueComponent) {

	var params []segmentSelectorParam
	var value valueComponent

	elems := strings.Split(s, "+")
	for i, e := range elems {

		comps := strings.Split(e, ":")
		for j, c := range comps {
			if len(c) == 0 {
				continue
			}
			if c == "?" {
				value = valueComponent{
					comp: i + 1,
					elem: j,
				}
				continue
			}

			params = append(params, segmentSelectorParam{
				comp:  i + 1,
				elem:  j,
				value: c,
			})
		}
	}

	return params, value
}

// selectValue returns a seqment element value.
func (sel segmentSelector) selectValue(seg spec.Segment) string {
	if sel.value.comp == 0 {
		return string(seg)
	}
	return seg.Comp(sel.value.comp).Elem(sel.value.elem)
}

// matches returns true if segment group path, segment tag and segment selector parameters matches.
func (sel segmentSelector) matches(node *spec.Node, seg spec.Segment) bool {

	if sel.tag != seg.Tag() {
		return false
	}
	if sel.path != node.Path() {
		return false
	}

	matches := true
	for _, param := range sel.params {
		matches = matches && param.matches(seg)
	}
	return matches
}

// matches returns true if a segment element value matches the param.
func (param segmentSelectorParam) matches(seg spec.Segment) bool {
	comp := seg.Comp(param.comp).Elem(param.elem)
	candidates := strings.Split(param.value, "|")
	for _, c := range candidates {
		if comp == c {
			return true
		}
	}
	return false
}
