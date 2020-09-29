package edifact

import "io"

// Unmarshaller interface
type Unmarshaller interface {
	UnmarshalEdifact(data []byte) error
}

// Unmarshal edifact document into data structure.
func Unmarshal(r io.Reader, target interface{}) error {
	return nil
}
