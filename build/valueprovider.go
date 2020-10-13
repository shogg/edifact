package build

import (
	"reflect"
)

// valueProvider interface for creating and connecting values.
type valueProvider interface {
	clearValue()
	setValue(reflect.Value)
	getValue(*decodeNode) reflect.Value
	newValue(*decodeNode) reflect.Value
}

type ptrProvider struct {
	ptr reflect.Value
}

func (p *ptrProvider) clearValue()              { p.ptr = reflect.Value{} }
func (p *ptrProvider) setValue(v reflect.Value) { p.ptr = v }
func (p *ptrProvider) getValue(node *decodeNode) reflect.Value {
	if !p.ptr.IsValid() {
		p.ptr = node.parent.getValue()
	}
	e := p.ptr.Elem()
	if !e.IsValid() {
		e = p.newValue(node)
	}
	return e
}
func (p *ptrProvider) newValue(node *decodeNode) reflect.Value {
	if !p.ptr.IsValid() {
		p.ptr = node.parent.getValue()
	}
	p.ptr.Set(reflect.New(p.ptr.Type().Elem()))
	return p.ptr.Elem()
}

type arrayProvider struct {
	array reflect.Value
	index int
}

func (p *arrayProvider) clearValue()              { p.array = reflect.Value{} }
func (p *arrayProvider) setValue(v reflect.Value) { p.array = v }
func (p *arrayProvider) getValue(node *decodeNode) reflect.Value {
	if !p.array.IsValid() {
		p.array = node.parent.getValue()
	}
	if p.index >= p.array.Len() {
		return reflect.Value{}
	}
	v := p.array.Index(p.index)
	return v
}
func (p *arrayProvider) newValue(node *decodeNode) reflect.Value {
	if !p.array.IsValid() {
		p.array = node.parent.getValue()
	}
	if p.index >= p.array.Len() {
		return reflect.Value{}
	}
	v := p.array.Index(p.index)
	p.index++
	return v
}

type sliceProvider struct {
	slice reflect.Value
}

func (p *sliceProvider) clearValue()              { p.slice = reflect.Value{} }
func (p *sliceProvider) setValue(v reflect.Value) { p.slice = v }
func (p *sliceProvider) getValue(node *decodeNode) reflect.Value {
	if !p.slice.IsValid() {
		p.slice = node.parent.getValue()
	}
	return p.slice.Index(p.slice.Len() - 1)
}
func (p *sliceProvider) newValue(node *decodeNode) reflect.Value {
	if !p.slice.IsValid() {
		p.slice = node.parent.getValue()
	}
	p.slice.Set(reflect.Append(p.slice, reflect.Zero(p.slice.Type().Elem())))

	return p.slice.Index(p.slice.Len() - 1)
}

type structProvider struct {
	str   reflect.Value
	field int
}

func (p *structProvider) clearValue()              { p.str = reflect.Value{} }
func (p *structProvider) setValue(v reflect.Value) { p.str = v }
func (p *structProvider) getValue(node *decodeNode) reflect.Value {
	if !p.str.IsValid() {
		p.str = node.parent.getValue()
	}
	return p.str.Field(p.field)
}
func (p *structProvider) newValue(node *decodeNode) reflect.Value {
	return p.getValue(node)
}

type structFieldProvider struct {
	field int
}

func (p *structFieldProvider) clearValue()              {}
func (p *structFieldProvider) setValue(v reflect.Value) {}
func (p *structFieldProvider) getValue(node *decodeNode) reflect.Value {
	node.parent.valProvider.(*structProvider).field = p.field
	return node.parent.getValue()
}
func (p *structFieldProvider) newValue(node *decodeNode) reflect.Value {
	return p.getValue(node)
}
