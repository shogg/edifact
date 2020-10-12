package spec

import (
	"fmt"
	"strings"
)

// Node part of message format spec
type Node struct {
	Tag          string
	Type         NodeType
	Required     Required
	Max          int // currently ignored
	Parent       *Node
	FirstChild   *Node
	Sibling      *Node
	SegmentGroup *Node

	path *string
}

// NodeType node type
type NodeType int

// node types
const (
	NodeMessage NodeType = iota
	NodeSegment
	NodeSegmentGroup
)

// Required mandatory, conditional
type Required int

// Madandatory, Conditional
const (
	M Required = iota
	C
)

// Msg creates a message node.
func Msg(tag string, children ...*Node) *Node {
	return newNode(NodeMessage, tag, C, 1, children)
}

// S creates a segment node.
func S(tag string, req Required, max int) *Node {
	return newNode(NodeSegment, tag, req, max, nil)
}

// SG creates a segment group node.
func SG(tag string, req Required, max int, children ...*Node) *Node {
	return newNode(NodeSegmentGroup, tag, req, max, children)
}

func newNode(nodeType NodeType, tag string, p Required, max int, children []*Node) *Node {

	n := &Node{
		Type:         nodeType,
		Tag:          tag,
		Required:     p,
		Max:          max,
		Parent:       nil,
		FirstChild:   nil,
		Sibling:      nil,
		SegmentGroup: nil,
	}

	lenChildren := len(children)
	if len(children) > 0 {
		n.FirstChild = children[0]
	}

	for i, c := range children {
		c.Parent = n
		if i+1 < lenChildren {
			c.Sibling = children[i+1]
		}
		if n.Type == NodeSegmentGroup {
			c.SegmentGroup = n
		}
	}

	return n
}

func (node *Node) FindNode(path, tag string) *Node {
	var s *Node
	node.iterate(true, func(n *Node) bool {
		if n.Type == NodeSegment && n.Tag == tag && n.Path() == path {
			s = n
			return false
		}
		return true
	})
	return s
}

func (node *Node) Path() string {
	if node.path != nil {
		return *node.path
	}

	var buf strings.Builder
	for _, sg := range node.SegmentGroups() {
		buf.WriteString(sg.Tag)
		buf.WriteByte('/')
	}
	node.path = new(string)
	*node.path = buf.String()
	return *node.path
}

// SegmentGroups all segment groups.
func (node *Node) SegmentGroups() []*Node {

	if node == nil {
		return nil
	}

	count := 0
	sg := node.SegmentGroup
	for sg != nil {
		count++
		sg = sg.SegmentGroup
	}

	result := make([]*Node, count)
	sg = node.SegmentGroup
	i := count - 1
	for sg != nil {
		result[i] = sg
		sg = sg.SegmentGroup
		i--
	}

	return result
}

// Transition switch to next node with matching segment tag.
func (node *Node) Transition(tag string) (*Node, error) {

	if node.Tag == tag {
		return node, nil
	}

	var s *Node
	node.iterate(false, func(n *Node) bool {
		if n.Type == NodeSegment && n.Tag == tag {
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
func (node *Node) iterate(all bool, f func(*Node) bool) {

	n := node
	for n != nil {

		if n.FirstChild != nil {
			n = n.FirstChild
		} else if n.Sibling != nil {
			n = n.Sibling
		} else {
			n = n.parentSibling(true, all)
			if n == nil {
				return
			}
		}

		if n != nil && !f(n) {
			return
		}

		if n.Required == M && !all {
			n = n.parentSibling(false, false)
			if n != nil && !f(n) {
				return
			}
		}
	}
}

func (node *Node) parentSibling(loop, loopOverride bool) *Node {
	n := node
	for n.Parent != nil {
		n = n.Parent
		if loop && !loopOverride && n.Type == NodeSegmentGroup {
			return n
		}
		if n.Sibling != nil {
			return n.Sibling
		}
		if !loopOverride {
			loop = true
		}
	}
	return nil
}
