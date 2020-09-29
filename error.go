package edifact

import (
	"fmt"
)

// Error edifact error
type Error int

// error definitions
const (
	ErrMissingSegmentDelimiter Error = iota
)

func (e Error) Error() string {
	switch e {
	case ErrMissingSegmentDelimiter:
		return "ErrMissingSegmentDelimiter"
	default:
		return fmt.Sprintf("parser.Error(%d)", e)
	}
}
