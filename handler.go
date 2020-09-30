package edifact

import "github.com/shogg/edifact/spec"

// Handler interface
type Handler interface {
	Handle(node *spec.Node, seg Segment) error
}
