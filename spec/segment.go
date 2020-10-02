package spec

import (
	"regexp"
	"strings"
)

type (
	// Segment is a sequence of elements joined by + ending in '
	Segment string
	// Element is a sequence of components joined by :
	Element string
	// Component string
)

var (
	regexTag             = regexp.MustCompile(`^[^+]*`)
	regexSplitElements   = regexp.MustCompile(`\+`)
	regexSplitComponents = regexp.MustCompile(`\:|'`)
)

// Tag retrieves the segment tag.
func (s Segment) Tag() string {
	return regexTag.FindString(string(s))
}

// Elem retrieves the ith element.
func (s Segment) Elem(i int) Element {
	elems := regexSplitElements.Split(string(s), -1)
	elems = concatAtReleaseChar(elems)
	if i >= len(elems) {
		return ""
	}
	return Element(elems[i])
}

// Comp retrieves the ith component.
func (e Element) Comp(i int) string {
	comps := regexSplitComponents.Split(string(e), -1)
	comps = concatAtReleaseChar(comps)
	if i >= len(comps) {
		return ""
	}
	return comps[i]
}

func concatAtReleaseChar(list []string) []string {

	releaseChar := false
	for i := range list {
		if len(list[i]) == 0 {
			continue
		}
		if list[i][len(list[i])-1] == '?' {
			releaseChar = true
		}
	}
	if !releaseChar {
		return list
	}

	result := make([]string, 0, len(list))
	for i := range list {
		var buf strings.Builder
		buf.WriteString(list[i])
		for j := i; j < len(list)-1; j++ {
			if list[j][len(list[j])-1] != '?' {
				break
			}
			buf.WriteByte('+')
			buf.WriteString(list[j+1])
		}
		result = append(result, buf.String())
	}

	return result
}
