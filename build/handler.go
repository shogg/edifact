package build

import (
	"github.com/shogg/edifact/spec"
)

// Handler edifact.Handler implementation to build arbitrary data structures.
type Handler struct {
	target     interface{}
	decodeTree decodeTree
}

// NewHandler creates a new handler.
func NewHandler(target interface{}) *Handler {
	return &Handler{target: target}
}

// Handle edifact.Handler implementation.
func (h *Handler) Handle(specNode *spec.Node, seg spec.Segment, loop bool) error {

	if h.decodeTree == nil {
		var err error
		h.decodeTree, err = newDecodeTree(specNode, h.target)
		if err != nil {
			return err
		}
	}

	decodeNodes := h.decodeTree[specNode.Key()]

	if specNode.Tag == "UNH" {
		for _, n := range decodeNodes {
			n.newValue()
		}
	}

	if loop {
		parentNodes := h.decodeTree[specNode.Parent.Key()]
		for _, p := range parentNodes {
			for _, c := range p.children {
				c.newValue()
			}
		}
	}

	for _, n := range decodeNodes {
		if err := n.decode(specNode.Path(), seg); err != nil {
			return err
		}
	}

	return nil
}
