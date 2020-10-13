package parse

import "github.com/shogg/edifact/spec"

// Handler parse handler interface
type Handler interface {
	Handle(*spec.Node, spec.Segment, bool) error
}
