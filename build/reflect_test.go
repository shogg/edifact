package build_test

import (
	"reflect"
	"testing"
)

func TestValue(t *testing.T) {

	var i int

	v := reflect.ValueOf(&i).Elem()
	v.SetInt(1)
	if i != 1 {
		t.Error("i == 1 exprected")
	}
}

func TestPointer(t *testing.T) {

	var i *int

	v := reflect.ValueOf(&i).Elem()
	v.Set(reflect.New(v.Type().Elem()))
	v.Elem().SetInt(1)
	if *i != 1 {
		t.Error("*i == 1 exprected")
	}
}

func TestArray(t *testing.T) {

	var a [3]int

	v := reflect.ValueOf(&a).Elem()
	v.Index(2).SetInt(1)
	if a[2] != 1 {
		t.Error("a[2] == 1 exprected")
	}
}

func TestAppend(t *testing.T) {

	var a []int

	v := reflect.ValueOf(&a).Elem()
	v.Set(reflect.Append(v, reflect.ValueOf(1)))
	if a[0] != 1 {
		t.Error("a[0] == 1 exprected")
	}
}

func TestAppendStruct(t *testing.T) {

	var a []struct {
		I int
	}

	v := reflect.ValueOf(&a).Elem()
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	v.Set(reflect.Append(v, reflect.Zero(v.Type().Elem())))
	v.Index(0).Field(0).SetInt(1)
	v.Index(1).Field(0).SetInt(2)
	if a[0].I != 1 {
		t.Error("a[0] == 1 exprected")
	}
	if a[1].I != 2 {
		t.Error("a[1] == 2 exprected")
	}
}

func TestStruct(t *testing.T) {

	type S struct {
		I int
	}

	var s S

	v := reflect.ValueOf(&s).Elem()
	v.Field(0).SetInt(1)
	if s.I != 1 {
		t.Error("s.I == 1 exprected")
	}
}

func TestEmptyInterface(t *testing.T) {

	type S struct {
		I int
	}

	var s S
	var i interface{} = &s

	v := reflect.ValueOf(i).Elem()
	v.Field(0).SetInt(1)
	if s.I != 1 {
		t.Error("s.I == 1 exprected")
	}
}

func TestStructStruct(t *testing.T) {

	type S struct {
		S struct {
			I int
		}
	}

	var s S

	v := reflect.ValueOf(&s).Elem()
	v.Field(0).Field(0).SetInt(1)
	if s.S.I != 1 {
		t.Error("s.S.I == 1 exprected")
	}
}

func TestStructStructPointer(t *testing.T) {

	type S struct {
		S *struct {
			I int
		}
	}

	var s S

	v := reflect.ValueOf(&s).Elem()
	v.Field(0).Set(reflect.New(v.Field(0).Type().Elem()))
	v.Field(0).Elem().Field(0).SetInt(1)
	if s.S.I != 1 {
		t.Error("s.S.I == 1 exprected")
	}
}
