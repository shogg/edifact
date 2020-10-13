package build

import (
	"fmt"
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/shogg/edifact/spec"
)

type decodeTree map[*spec.Node][]*decodeNode

type decodeNode struct {
	parent      *decodeNode
	children    []*decodeNode
	valProvider valueProvider
	segSelector segmentSelector
	specNode    *spec.Node
}

func newDecodeTree(specNode *spec.Node, target interface{}) (decodeTree, error) {

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return nil, fmt.Errorf("target is nil or not a pointer: %s", v)
	}

	tree := decodeTree{}
	v = v.Elem()
	root := addDecodeNode(tree, nil, specNode, v.Type())
	root.setValue(v)
	tree.add(specNode, root)

	return tree, nil
}

func (tree decodeTree) add(specNode *spec.Node, node *decodeNode) {
	if specNode == nil {
		return
	}
	list := tree[specNode]
	list = append(list, node)
	tree[specNode] = list
}

func (node *decodeNode) setValue(v reflect.Value) {
	node.valProvider.setValue(v)
}

func (node *decodeNode) getValue() reflect.Value {
	v := node.valProvider.getValue(node)
	return v
}

func (node *decodeNode) newValue() reflect.Value {
	v := node.valProvider.newValue(node)
	for _, c := range node.children {
		c.clearValue()
	}
	return v
}

func (node *decodeNode) clearValue() {
	node.valProvider.clearValue()
	for _, c := range node.children {
		c.clearValue()
	}
}

func addDecodeNode(
	tree decodeTree,
	parent *decodeNode,
	specNode *spec.Node,
	t reflect.Type) *decodeNode {

	node := &decodeNode{parent: parent}

	switch t.Kind() {
	case reflect.Ptr:
		node.valProvider = &ptrProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Array:
		node.valProvider = &arrayProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Slice:
		node.valProvider = &sliceProvider{}
		addDecodeNode(tree, node, specNode, t.Elem())
	case reflect.Struct:
		node.valProvider = &structProvider{}
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			r, _ := utf8.DecodeRuneInString(f.Name)
			if !unicode.IsUpper(r) {
				continue
			}
			sel := parseSegmentSelector(f.Tag.Get("edifact"))
			sn := specNode.FindNode(sel.path, sel.tag)
			n := &decodeNode{
				parent:      node,
				segSelector: sel,
				valProvider: &structFieldProvider{field: i},
				specNode:    sn,
			}
			node.children = append(node.children, n)

			tree.add(sn, n)

			addDecodeNode(tree, n, specNode, f.Type)
		}
	default:
		return nil
	}

	if parent != nil {
		parent.children = append(parent.children, node)
	}
	return node
}

func (node *decodeNode) decode(seg spec.Segment) error {
	if !node.segSelector.matches(node.specNode, seg) {
		return nil
	}

	s := node.segSelector.selectValue(seg)
	v := node.getValue()
	return decode(s, v)
}
