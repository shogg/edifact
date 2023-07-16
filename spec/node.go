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

// Path segment group path as string.
func (node *Node) Path() string {
	if node.path != nil {
		return *node.path
	}

	var buf strings.Builder
	for _, sg := range node.SegmentGroups() {
		buf.WriteString(sg.Tag)
		buf.WriteByte('/')
	}
	path := buf.String()
	node.path = &path
	return *node.path
}

// Key segment group path and segment tag.
func (node *Node) Key() string {
	var buf strings.Builder
	buf.WriteString(node.Path())
	buf.WriteByte('/')
	buf.WriteString(node.Tag)
	return buf.String()
}

// SegmentGroups segment group path.
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

// FindNode find a node with matching segment group path and segment tag.
func (node *Node) FindNode(path, tag string) *Node {

	if node.Tag == tag && node.Path() == path {
		return node
	}

	var s *Node
	node.findNode(func(n *Node) bool {
		if n.Tag == tag && n.Path() == path {
			s = n
			return true
		}
		return false
	})
	return s
}

// Iteration order: (1) firstChild (!leave), (2) sibling, (3) parent (leave)
//
//	-->o-->o
//	(1)|\
//	   v \(3)
//	   o-->o
//	   (2)
func (node *Node) findNode(f func(*Node) bool) {

	var leaveSG bool

	n := node
	for {
		if n.FirstChild != nil && !leaveSG {
			n = n.FirstChild
		} else if n.Sibling != nil {
			n = n.Sibling
			leaveSG = false
		} else if n.Parent != nil {
			n = n.Parent
			leaveSG = true
		} else {
			return
		}

		if f(n) {
			return
		}
	}
}

// Transition switch to next node with matching segment tag.
func (node *Node) Transition(tag string) (*Node, error) {

	if node.Tag == tag {
		return node, nil
	}

	var s *Node
	node.transition(func(n *Node) bool {
		if n.Tag == tag {
			s = n
			return true
		}
		return false
	})

	if s == nil {
		return nil, fmt.Errorf("%w: %s", ErrUnexpectedSegment, tag)
	}
	return s, nil
}

func (node *Node) transition(f func(*Node) bool) {

	var leaveSG bool

	n := node
	for {
		if n.FirstChild != nil && !leaveSG {
			n = n.FirstChild
		} else if n.Sibling != nil {
			n = n.Sibling
			leaveSG = false
		} else if n.Parent != nil {
			n = n.Parent
			leaveSG = false
		} else {
			return
		}

		if f(n) {
			return
		}

		if n.Required == M {
			n = n.Parent
			leaveSG = true
		}
	}
}
