package spec

// Node part of message format spec
type Node struct {
	Tag          string
	Type         Type
	Required     Required
	Max          int // currently ignored
	Parent       *Node
	FirstChild   *Node
	Sibling      *Node
	SegmentGroup *Node
	Level        int // present when traversed (edifact.parser.Transition)
}

// Type node type
type Type int

// node types
const (
	Message Type = iota
	Segment
	SegmentGroup
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
	return newNode(Message, tag, C, 1, children)
}

// S creates a segment node.
func S(tag string, req Required, max int) *Node {
	return newNode(Segment, tag, req, max, nil)
}

// SG creates a segment group node.
func SG(tag string, req Required, max int, children ...*Node) *Node {
	return newNode(SegmentGroup, tag, req, max, children)
}

func newNode(nodeType Type, tag string, p Required, max int, children []*Node) *Node {

	n := &Node{
		Type:         nodeType,
		Tag:          tag,
		Required:     p,
		Max:          max,
		Parent:       nil,
		FirstChild:   nil,
		Sibling:      nil,
		SegmentGroup: nil,
		Level:        0,
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
		if n.Type == SegmentGroup {
			c.SegmentGroup = n
		}
	}

	return n
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
