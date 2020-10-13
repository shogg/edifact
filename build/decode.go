package build

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var typeTime = reflect.TypeOf(time.Time{})

func decode(s string, v reflect.Value) error {

	// edifact.Unmarshaller
	unmarshal := v.MethodByName("UnmarshalEdifact")
	if unmarshal.IsValid() {
		in := []reflect.Value{reflect.ValueOf([]byte(s))}
		out := unmarshal.Call(in)
		if out[0].IsNil() {
			return nil
		}
		return out[0].Interface().(error)
	}

	// time.Time
	if v.Type().AssignableTo(typeTime) {
		return decodeTime(s, v)
	}

	// simple data types
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

func decodeTime(s string, v reflect.Value) error {

	if s[:3] == "DTM" {
		comp := strings.SplitN(s, ":", 3)
		s = comp[1]
	}

	t, err := time.Parse("20060102", s)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(t))
	return nil
}
