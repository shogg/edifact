package build

import "github.com/shogg/edifact/spec"

type Handler struct {
	Target interface{}
}

func (h *Handler) Handle(n *spec.Node, s spec.Segment) error {
	return nil
}
