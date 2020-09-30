package edifact

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/shogg/edifact/spec"
)

// parser parses edifact messages.
type parser struct {
	scanner   *bufio.Scanner
	node      *spec.Node
	lineNr    int
	segmentNr int
}

func newParser(r io.Reader) *parser {

	s := bufio.NewScanner(r)
	s.Split(segments('\''))

	return &parser{
		scanner: s,
	}
}

func (p *parser) parse(h Handler) error {

	for p.scanner.Scan() {
		seg := Segment(p.scanner.Text())

		p.segmentNr++
		if strings.ContainsAny(string(seg), "\r\n") {
			p.lineNr++
		}
		seg = Segment(strings.TrimSpace(string(seg)))

		if seg.Tag() == "UNH" {
			p.node = spec.Get(seg.Elem(2).Comp(0))
		}

		if p.node != nil {
			next, err := p.node.Transition(seg.Tag())
			if err != nil {
				return p.annotate(err)
			}
			p.node = next
		}

		if err := h.Handle(p.node, seg); err != nil {
			return p.annotate(err)
		}
	}

	return p.annotate(p.scanner.Err())
}

func (p *parser) annotate(err error) error {
	return fmt.Errorf("line %d, segment %d: %w", p.lineNr, p.segmentNr, err)
}

func segments(del byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		index := bytes.IndexByte(data, del)
		if index < 0 && !atEOF {
			if len(bytes.TrimSpace(data)) == 0 {
				return len(data), nil, nil
			}
			return 0, nil, ErrMissingSegmentDelimiter
		}
		return index + 1, data[:index+1], nil
	}
}
