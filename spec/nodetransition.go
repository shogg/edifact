package spec

import (
	"fmt"
)

// Transition switch to next node with matching segment tag.
func (node *Node) Transition(tag string) (*Node, error) {

	if node.Tag == tag {
		return node, nil
	}

	var s *Node
	node.iterate(func(n *Node) bool {
		if n.Type == Segment && n.Tag == tag {
			s = n
			return false
		}
		return true
	})

	if s == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedSegment, tag)
	}
	return s, nil
}

// Iteration order: (1) child, (2) sibling, (3) parent sibling.
// If parent is a group repeat the group.
// --->o -->o
//  (1)|  \
//     v   \(3)
// 	   o--->o
//     (2)
func (node *Node) iterate(f func(*Node) bool) {

	n := node
	for n != nil {

		if n.FirstChild != nil {
			n.FirstChild.Level = n.Level + 1
			n = n.FirstChild
		} else if n.Sibling != nil {
			n.Sibling.Level = n.Level
			n = n.Sibling
		} else {
			n = n.parentSibling(true)
			if n == nil {
				return
			}
		}

		if !f(n) {
			return
		}

		if n.Required == M {
			n = n.parentSibling(false)
			if !f(n) {
				return
			}
		}
	}
}

func (node *Node) parentSibling(loop bool) *Node {
	n := node
	level := node.Level
	for n.Parent != nil {
		n = n.Parent
		if loop && n.Type == SegmentGroup {
			return n
		}
		level--
		if n.Sibling != nil {
			n.Sibling.Level = level
			return n.Sibling
		}
	}
	return nil
}
