package build

import "fmt"

// Error build error
type Error int

// error definitions
const (
	ErrNotPointer Error = iota
	ErrNotStruct
	ErrNotImplemented
)

func (e Error) Error() string {
	if msg, ok := errMessages[e]; ok {
		return msg
	}
	return fmt.Sprintf("build.Error(%d)", e)
}

var errMessages = map[Error]string{
	ErrNotPointer:     "not pointer",
	ErrNotStruct:      "not struct",
	ErrNotImplemented: "not implemented",
}
