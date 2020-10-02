package build

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/shogg/edifact/spec"
)

func unmarshal(v interface{}, seg spec.Segment, comp ValueComponent) error {
	return unmarshalStruct(reflect.ValueOf(v), seg, comp)
}

func unmarshalStruct(p reflect.Value, seg spec.Segment, comp ValueComponent) error {

	if p.Kind() != reflect.Ptr || p.IsNil() {
		return fmt.Errorf("%w: %s", ErrNotPointer, p.Type())
	}

	s := p.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("%w: %s", ErrNotStruct, s.Type())
	}

	for i := 0; i < s.NumField(); i++ {
		f := s.Field(i)
		unmarshalLiteral(f, seg, comp)
	}

	return nil
}

func unmarshalLiteral(v reflect.Value, seg spec.Segment, comp ValueComponent) error {

	s := seg.Elem(comp.Elem).Comp(comp.Comp)

	switch v.Kind() {
	case reflect.String:
		v.SetString(s)
	case reflect.Bool:
		if s == "" {
			v.SetBool(false)
			return nil
		}
		b, err := strconv.ParseBool(s)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Complex64, reflect.Complex128:
		n, err := strconv.ParseComplex(s, 128)
		if err != nil {
			return err
		}
		v.SetComplex(n)
	default:
		return fmt.Errorf("%s: %s", ErrNotImplemented, v.Type())
	}

	return nil
}
