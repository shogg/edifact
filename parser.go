package edifact

import (
	"bufio"
	"bytes"
	"io"

	"github.com/shogg/edifact/internal/sgstack"
	"github.com/shogg/edifact/spec"
)

// parser parses edifact messages.
type parser struct {
	scanner *bufio.Scanner
	state   state
}

// state holds the current parser state.
type state struct {
	node          *spec.Node
	segmentNr     int
	segmentGroups []*spec.Node
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
		if seg.Tag() == "UNH" {
			p.state.node = spec.Get(seg.Elem(2).Comp(0))
		}

		if p.state.node != nil {
			next, err := p.state.node.Transition(seg.Tag())
			if err != nil {
				return err
			}
			if p.state.node.Level < next.Level {
				sgstack.Push(&p.state.segmentGroups, next.SegmentGroup)
			} else if p.state.node.Level > next.Level {
				for sgstack.Pop(&p.state.segmentGroups).Tag != next.SegmentGroup.Tag {
				}
			}

			p.state.node = next
		}

		if err := h.Handle(p.state.segmentGroups, seg); err != nil {
			return err
		}
	}

	return p.scanner.Err()
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
		token = bytes.TrimLeft(data[:index+1], "\r\n\t ")
		return index + 1, token, nil
	}
}
