package edifact

import (
	"io"

	"github.com/shogg/edifact/internal/build"
	"github.com/shogg/edifact/parse"
)

// Unmarshaller interface for custom data types.
type Unmarshaller interface {
	UnmarshalEdifact(data []byte) error
}

// Unmarshal edifact document into target data structure.
func Unmarshal(r io.Reader, target interface{}) error {
	h := &build.Handler{Target: target}
	return parse.Parse(r, h)
}
