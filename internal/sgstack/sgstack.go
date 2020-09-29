package sgstack

import "github.com/shogg/edifact/spec"

// Pop retrieve the top node and drop it.
func Pop(st *[]*spec.Node) *spec.Node {
	top := (*st)[len(*st)-1]
	*st = (*st)[:len(*st)-1]
	return top
}

// Push append a new top node.
func Push(st *[]*spec.Node, sg *spec.Node) {
	*st = append(*st, sg)
}
