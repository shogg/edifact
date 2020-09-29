package spec

// Transition switch to next node with matching segment tag.
// If the node is in a different segment group
// switch to the node of the segment group.
func (node *Node) Transition(tag string) (*Node, error) {

	if node.Tag == tag {
		return node, nil
	}

	var sg *Node
	var s *Node
	node.iterate(func(n *Node) bool {
		switch n.Type {
		case SegmentGroup:
			sg = n
		case Segment:
			if n.Tag == tag {
				s = n
				return false
			}
		}
		return true
	})

	if sg != nil && s != nil {
		return sg, nil
	}
	if s != nil {
		return s, nil
	}
	return nil, ErrUndefinedSegment
}

// iteration order: (1) child, (2) sibling, (3) parent sibling. Recursion in this order.
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
			n = n.parentSibling()
			if n == nil {
				return
			}
		}

		if !f(n) {
			return
		}

		if n.Required == M {
			n = n.parentSibling()
		}
	}
}

func (node *Node) parentSibling() *Node {
	n := node
	level := node.Level
	for n.Parent != nil {
		n = n.Parent
		level--
		if n.Sibling != nil {
			n.Sibling.Level = level
			return n.Sibling
		}
	}
	return nil
}
