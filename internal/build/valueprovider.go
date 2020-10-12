package build

import (
	"reflect"
)

type ValueProvider interface {
	SetValue(v reflect.Value)
	NewValue(node *DecodeNode) reflect.Value
}

type PtrProvider struct {
	ptr reflect.Value
}

func (p *PtrProvider) SetValue(v reflect.Value) { p.ptr = v }
func (p *PtrProvider) NewValue(node *DecodeNode) reflect.Value {
	if !p.ptr.IsValid() {
		p.ptr = node.Parent.NewValue()
	}
	p.ptr.Set(reflect.New(p.ptr.Type().Elem()))
	return p.ptr.Elem()
}

type ArrayProvider struct {
	array reflect.Value
	index int
}

func (p *ArrayProvider) SetValue(v reflect.Value) { p.array = v }
func (p *ArrayProvider) NewValue(node *DecodeNode) reflect.Value {
	if !p.array.IsValid() {
		p.array = node.Parent.NewValue()
	}
	if p.index >= p.array.Len() {
		return reflect.Value{}
	}
	v := p.array.Index(p.index)
	p.index++
	return v
}

type SliceProvider struct {
	slice reflect.Value
}

func (p *SliceProvider) SetValue(v reflect.Value) { p.slice = v }
func (p *SliceProvider) NewValue(node *DecodeNode) reflect.Value {
	if !p.slice.IsValid() {
		p.slice = node.Parent.NewValue()
	}
	p.slice.Set(reflect.Append(p.slice, reflect.Zero(p.slice.Type().Elem())))

	return p.slice.Index(p.slice.Len() - 1)
}

type StructProvider struct {
	str   reflect.Value
	field int
}

func (p *StructProvider) SetValue(v reflect.Value) { p.str = v }
func (p *StructProvider) NewValue(node *DecodeNode) reflect.Value {
	if !p.str.IsValid() {
		p.str = node.Parent.NewValue()
	}

	return p.str.Field(p.field)
}

type StructFieldProvider struct {
	field int
}

func (p *StructFieldProvider) SetValue(v reflect.Value) {}
func (p *StructFieldProvider) NewValue(node *DecodeNode) reflect.Value {
	node.Parent.ValueProvider.(*StructProvider).field = p.field
	return node.Parent.NewValue()
}
