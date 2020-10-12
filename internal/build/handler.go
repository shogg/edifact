package build

import (
	"github.com/shogg/edifact/spec"
)

type Handler struct {
	Target     interface{}
	decodeTree DecodeTree
}

func (h *Handler) Handle(specNode *spec.Node, seg spec.Segment) error {

	if h.decodeTree == nil {
		var err error
		h.decodeTree, err = newDecodeTree(specNode, h.Target)
		if err != nil {
			return err
		}
	}

	decodeNodes := h.decodeTree[specNode]
	for _, n := range decodeNodes {
		if err := n.Decode(seg); err != nil {
			return err
		}
	}

	return nil
}
