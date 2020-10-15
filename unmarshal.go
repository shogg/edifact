package edifact

import (
	"io"

	"github.com/shogg/edifact/build"
	"github.com/shogg/edifact/parse"
)

// Unmarshaller interface for custom data type parsing.
type Unmarshaller interface {
	UnmarshalEdifact(data []byte) error
}

// Unmarshal edifact document into target data structure.
func Unmarshal(r io.Reader, target interface{}) error {
	h := build.NewHandler(target)
	return parse.Parse(r, h)
}
