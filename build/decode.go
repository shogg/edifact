package build

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type unmarshaller interface {
	UnmarshalEdifact(data []byte) error
}

var (
	typeTime = reflect.TypeOf(time.Time{})
)

func decode(s string, v reflect.Value) error {

	// edifact.Unmarshaller
	p := v
	if p.Kind() != reflect.Ptr {
		p = v.Addr()
	}
	if iface, ok := p.Interface().(unmarshaller); ok {
		if p.IsNil() {
			p.Set(reflect.New(p.Type().Elem()))
		}
		return iface.UnmarshalEdifact([]byte(s))
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
		if s == "" {
			v.SetInt(0)
			return nil
		}
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(n)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		if s == "" {
			v.SetUint(0)
			return nil
		}
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(n)
	case reflect.Float32, reflect.Float64:
		if s == "" {
			v.SetFloat(0.0)
			return nil
		}
		n, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}
		v.SetFloat(n)
	case reflect.Complex64, reflect.Complex128:
		if s == "" {
			v.SetComplex(0)
			return nil
		}
		n, err := strconv.ParseComplex(s, 128)
		if err != nil {
			return err
		}
		v.SetComplex(n)
	case reflect.Slice:
		zeroValue := reflect.Zero(v.Type().Elem())
		v.Set(reflect.Append(v, zeroValue))
		lastEntry := v.Index(v.Len() - 1)
		return decode(s, lastEntry)
	default:
		return fmt.Errorf("%w: %s", ErrNotImplemented, v.Type())
	}

	return nil
}

const (
	dtmFormat102 = 102
	dtmFormat201 = 201
	dtmFormat203 = 203
)

func decodeTime(s string, v reflect.Value) error {

	dtmFormat := dtmFormat102
	if s[:3] == "DTM" {
		s = strings.Trim(s, "'")
		comp := strings.SplitN(s, ":", 3)
		s = comp[1]
		dtmFormat, _ = strconv.Atoi(comp[2])
	}

	format := "20060102"
	switch dtmFormat {
	case dtmFormat102:
		format = "20060102"
	case dtmFormat201, dtmFormat203:
		format = "200601021504"
	}

	t, err := time.Parse(format, s)
	if err != nil {
		return err
	}

	v.Set(reflect.ValueOf(t))
	return nil
}
