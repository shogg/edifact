package spec

import (
	"fmt"
)

// Error spec error
type Error int

// errors
const (
	ErrUndefinedSegment Error = iota
)

func (e Error) Error() string {
	switch e {
	case ErrUndefinedSegment:
		return "ErrUndefinedSegment"
	default:
		return fmt.Sprintf("spec.Error(%d)", e)
	}
}
