package spec

import (
	"fmt"
)

// Error spec error
type Error int

// errors
const (
	ErrUnexpectedSegment Error = iota
)

func (e Error) Error() string {
	if msg, ok := errMessages[e]; ok {
		return msg
	}
	return fmt.Sprintf("spec.Error(%d)", e)
}

var errMessages = map[Error]string{
	ErrUnexpectedSegment: "unexpected segment",
}
