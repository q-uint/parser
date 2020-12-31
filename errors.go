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

// ExpectError is an error that occurs on when an invalid/unsupported value is
// passed to the Parser.Expect function.
type ExpectError struct {
	Message string
}

func (e *ExpectError) Error() string {
	return fmt.Sprintf("expect: %s", e.Message)
}

// ExpectedParseError indicates that the parser Expected a different value than
// the Actual value present in the buffer.
type ExpectedParseError struct {
	// The value that was expected.
	Expected interface{}
	// The value it actually got.
	Actual string
}

func (e *ExpectedParseError) Error() string {
	return fmt.Sprintf("parse: expected %T %v but got %q", e.Expected, e.Expected, e.Actual)
}
