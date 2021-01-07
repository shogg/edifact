package parse

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/shogg/edifact/spec"
)

// parser annotates parse errors with line and segment number.
type parser struct {
	lineNr    int
	segmentNr int
}

func (p *parser) annotate(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("line %d, segment %d: %w", p.lineNr, p.segmentNr, err)
}

// Parse parses an edifact document.
func Parse(r io.Reader, h Handler) error {

	s := bufio.NewScanner(r)
	s.Split(segments('\'', '?'))
	p := &parser{lineNr: 1}

	var node *spec.Node

	for s.Scan() {
		token := s.Text()

		p.lineNr += strings.Count(token, "\n")
		p.segmentNr++

		token = strings.TrimSpace(token)
		if len(token) == 0 {
			break
		}
		if token[len(token)-1] != '\'' {
			return p.annotate(ErrMissingSegmentTerminator)
		}

		seg := spec.Segment(token)

		if seg.Tag() == "UNH" {
			msgType := seg.Comp(2).Elem(0)
			node = spec.Get(msgType)
			if node == nil {
				return p.annotate(fmt.Errorf(
					"unknown edifact message type: %s",
					msgType))
			}
		}

		if node != nil {
			next, loop, err := node.Transition(seg.Tag())
			if err != nil {
				return p.annotate(err)
			}
			node = next

			if err := h.Handle(node, seg, loop); err != nil {
				return p.annotate(err)
			}
		}
	}

	return p.annotate(s.Err())
}

func segments(delimiter, release byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		d := data
		index := 0
		for {
			i := bytes.IndexByte(d, delimiter)
			if i < 0 {
				if atEOF {
					return len(data), data, nil
				}
				return 0, nil, nil
			}
			index += i
			if !isReleased(d, i, release) {
				break
			}
			index++
			d = d[i+1:]
		}

		return index + 1, data[:index+1], nil
	}
}

func isReleased(data []byte, index int, release byte) bool {

	released := false
	for i := index - 1; i >= 0 && data[i] == release; i-- {
		released = !released
	}

	return released
}
