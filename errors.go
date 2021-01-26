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

// ExpectedParseError creates an ExpectedParseError error based on the given
// start and end cursor. Resets the parser tot the start cursor.
func (p *Parser) ExpectedParseError(expected interface{}, start, end *Cursor) *ExpectedParseError {
	if end == nil {
		end = start
	}
	defer p.Jump(start)
	return &ExpectedParseError{
		Expected: expected,
		String:   p.Slice(start, end),
		Conflict: *end,
	}
}

// ExpectedParseError indicates that the parser Expected a different value than
// the Actual value present in the buffer.
type ExpectedParseError struct {
	// The value that was expected.
	Expected interface{}
	// The value it actually got.
	String string
	// The position of the conflicting value.
	Conflict Cursor
}

func (e *ExpectedParseError) Error() string {
	var expected string
	switch v := e.Expected.(type) {
	case rune:
		expected = fmt.Sprintf("'%s'", string(v))
	case string:
		expected = fmt.Sprintf("%q", v)
	default:
		expected = fmt.Sprintf("%v", v)
	}

	got := e.String
	if len(e.String) == 1 {
		got = fmt.Sprintf("'%s'", string([]rune(e.String)[0]))
	} else {
		got = fmt.Sprintf("%q", e.String)
	}

	return fmt.Sprintf(
		"parse conflict [%02d:%03d]: expected %T %s but got %s",
		e.Conflict.row, e.Conflict.column, e.Expected, expected, got,
	)
}

// UnsupportedType indicates the type of the value is unsupported.
type UnsupportedType struct {
	Value interface{}
}

func (e *UnsupportedType) Error() string {
	return fmt.Sprintf("parse: value of type %T are not supported", e.Value)
}
