package spec

import (
	"strings"
)

type (
	// Segment is a sequence of composites joined by + ending in '
	Segment string
	// Composite is a sequence of elements joined by :
	Composite string
	// Element string
)

// Tag retrieves the segment tag.
func (s Segment) Tag() string {
	i := strings.IndexByte(string(s), '+')
	if i < 0 {
		return string(s)
	}
	return string(s)[:i]
}

// Comp retrieves the ith composite.
func (s Segment) Comp(i int) Composite {
	scanner := newScanner(string(s), '+')

	j := 0
	for scanner.Scan() {
		if j == i {
			return Composite(scanner.Text())
		}
		j++
	}

	return ""
}

// Elem retrieves the ith element.
func (e Composite) Elem(i int) string {
	scanner := newScanner(string(e), ':')

	j := 0
	for scanner.Scan() {
		if j == i {
			return removeMetaChars.Replace(scanner.Text())
		}
		j++
	}

	return ""
}

var removeMetaChars = strings.NewReplacer(
	"?'", "'",
	"??", "?",
	"'", "",
	"?", "",
)

type segmentScanner struct {
	str  string
	del  byte
	text string
}

func newScanner(str string, del byte) *segmentScanner {
	return &segmentScanner{
		str: str,
		del: del,
	}
}

func (s *segmentScanner) Scan() bool {
	if len(s.str) == 0 {
		return false
	}

	tmp := s.str
	var index int
	for {
		i := strings.IndexByte(tmp, s.del)
		if i < 0 {
			index += len(tmp)
			break
		}
		index += i
		if !isReleased([]byte(tmp), index, '?') {
			break
		}
		index++
		tmp = tmp[i+1:]
	}

	s.text = s.str[:index]
	s.str = s.str[index:]

	// if present remove delimiter
	if len(s.str) > 0 {
		s.str = s.str[1:]
	}

	return true
}

func (s *segmentScanner) Text() string {
	return s.text
}

// isReleased checks if the character at index is released by
// a release (escape) character in front of it.
func isReleased(data []byte, index int, release byte) bool {

	released := false
	for i := index - 1; i >= 0 && data[i] == release; i-- {
		released = !released
	}

	return released
}
