package parser

import "fmt"

// InitError is an error that occurs on instantiating new structures.
type InitError struct {
	// The error message. This should provide an intuitive message or advice on
	// how to solve this error.
	Message string
}

func (e *InitError) Error() string {
	return fmt.Sprintf("parser: %s", e.Message)
}
