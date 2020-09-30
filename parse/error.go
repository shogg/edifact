package parse

import (
	"fmt"
)

// Error parse error
type Error int

// error definitions
const (
	ErrMissingSegmentTerminator Error = iota
)

func (e Error) Error() string {
	if msg, ok := errMessages[e]; ok {
		return msg
	}
	return fmt.Sprintf("spec.Error(%d)", e)
}

var errMessages = map[Error]string{
	ErrMissingSegmentTerminator: "missing segment terminator",
}
