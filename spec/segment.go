package spec

import (
	"bufio"
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
	scanner := bufio.NewScanner(strings.NewReader(string(s)))
	scanner.Buffer(make([]byte, 64), 1024)
	scanner.Split(delimiter('+'))

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
	scanner := bufio.NewScanner(strings.NewReader(string(e)))
	scanner.Buffer(make([]byte, 64), 1024)
	scanner.Split(delimiter(':'))

	j := 0
	for scanner.Scan() {
		if j == i {
			return scanner.Text()
		}
		j++
	}

	return ""
}

func delimiter(del byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		for i := range data {
			if data[i] != del {
				continue
			}
			if i > 0 && data[i-1] == '?' {
				continue
			}
			return i + 1, data[:i], nil
		}
		if len(data) > 0 && data[len(data)-1] == '\'' {
			return len(data), data[:len(data)-1], nil
		}
		return len(data), data, nil
	}
}
