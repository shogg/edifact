package parse

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

// Parse parses an edifact document.
func Parse(r io.Reader, h Handler) error {

	s := bufio.NewScanner(r)
	s.Split(segments('\''))
	p := &parser{scanner: s, lineNr: 1}

	for p.scanner.Scan() {
		token := p.scanner.Text()

		p.segmentNr++
		p.lineNr += strings.Count(token, "\n")

		token = strings.TrimSpace(token)
		if token == "" {
			continue
		}

		seg := spec.Segment(token)

		if seg.Tag() == "UNH" {
			p.node = spec.Get(seg.Elem(2).Comp(0))
		}

		if p.node != nil {
			next, loop, err := p.node.Transition(seg.Tag())
			if err != nil {
				return p.annotate(err)
			}
			p.node = next

			if err := h.Handle(p.node, seg, loop); err != nil {
				return p.annotate(err)
			}
		}
	}

	return p.annotate(p.scanner.Err())
}

func (p *parser) annotate(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("line %d, segment %d: %w", p.lineNr, p.segmentNr, err)
}

func segments(del byte) bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF {
			return 0, nil, nil
		}
		index := bytes.IndexByte(data, del)
		if index < 0 && !atEOF {
			return 0, nil, nil
		}
		return index + 1, data[:index+1], nil
	}
}
