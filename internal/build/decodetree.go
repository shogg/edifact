package build

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/shogg/edifact/spec"
)

type DecodeTree map[*spec.Node][]*DecodeNode

type DecodeNode struct {
	Parent          *DecodeNode
	ValueProvider   ValueProvider
	SegmentSelector SegmentSelector
	SpecNode        *spec.Node
}

func newDecodeTree(specNode *spec.Node, target interface{}) (DecodeTree, error) {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("target is nil or not a pointer: %s", v)
	}

	tree := DecodeTree{}
	v = v.Elem()
	root := addDecodeNode(tree, nil, specNode, v.Type())
	root.ValueProvider.SetValue(v)

	return tree, nil
}

func (tree DecodeTree) Add(specNode *spec.Node, node *DecodeNode) {
	if specNode == nil {
		return
	}
	list := tree[specNode]
	list = append(list, node)
	tree[specNode] = list
}

func (node *DecodeNode) NewValue() reflect.Value {
	v := node.ValueProvider.NewValue(node)
	return v
}

func addDecodeNode(
	tree DecodeTree,
	parent *DecodeNode,
	specNode *spec.Node,
	t reflect.Type) *DecodeNode {

	node := &DecodeNode{Parent: parent}

	switch t.Kind() {
	case reflect.Ptr:
		node.ValueProvider = &PtrProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Array:
		node.ValueProvider = &ArrayProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Slice:
		node.ValueProvider = &SliceProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Struct:
		node.ValueProvider = &StructProvider{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			r, _ := utf8.DecodeRuneInString(f.Name)
			if !unicode.IsUpper(r) {
				continue
			}
			sel := ParseSegmentSelector(f.Tag.Get("edifact"))
			sn := specNode.FindNode(sel.Path, sel.Tag)
			n := &DecodeNode{
				Parent:          node,
				SegmentSelector: sel,
				ValueProvider:   &StructFieldProvider{field: i},
				SpecNode:        sn,
			}

			tree.Add(sn, n)

			addDecodeNode(tree, n, specNode, f.Type)
		}
	default:
		return nil
	}

	return node
}

func (node *DecodeNode) Decode(seg spec.Segment) error {
	if !node.SegmentSelector.Matches(node.SpecNode, seg) {
		return nil
	}

	s := node.SegmentSelector.Select(seg)
	v := node.NewValue()
	return decode(s, v)
}
