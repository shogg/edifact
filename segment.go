package edifact

import (
	"regexp"
	"strings"
)

type (
	// Segment sequence of elements joined by + ending in '
	Segment string
	// Element sequence of components joined by :
	Element string
	// Component string value
	Component string
)

var (
	regexTag             = regexp.MustCompile(`^[^+]*`)
	regexSplitElements   = regexp.MustCompile(`\+`)
	regexSplitComponents = regexp.MustCompile(`\:`)
)

// Tag retrieves the segment tag.
func (s Segment) Tag() string {
	return regexTag.FindString(string(s))
}

// Elem retrieves the ith element.
func (s Segment) Elem(i int) Element {
	elems := regexSplitElements.Split(string(s), -1)
	elems = concatAtReleaseChar(elems)
	return Element(elems[i])
}

// Comp retrieves the ith component.
func (e Element) Comp(i int) string {
	comps := regexSplitComponents.Split(string(e), -1)
	comps = concatAtReleaseChar(comps)
	return comps[i]
}

func concatAtReleaseChar(list []string) []string {

	releaseChar := false
	for i := range list {
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
		for j := i; j < len(list); j++ {
			if list[j][len(list)-1] != '?' {
				break
			}
			buf.WriteString(list[j])
		}
		result = append(result, buf.String())
	}

	return result
}
