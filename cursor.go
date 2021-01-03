package parser

import "fmt"

// Cursor allows you to record your current position so you can return to it
// later. Keeps track of its own position in the buffer of the parser.
type Cursor struct {
	// Rune is the value that the cursor points at.
	Rune rune
	// The position in the buffer.
	position int
}

func (c *Cursor) String() string {
	return fmt.Sprintf("%U %c", c.Rune, c.Rune)
}

// state manages the state of a parser. It contains a pointer to the last
// successfully parsed rune.
type state struct {
	end *Cursor
	p   *Parser
}

// Ok updates the end of the state and jumps to the next rune.
func (s *state) Ok(last *Cursor) {
	if last == nil {
		// Optional values have no last mark.
		return
	}
	s.end = last
	// We jump to the given cursor (last parsed rune) because it is not
	// guaranteed that the already parser did not pass it.
	s.p.Jump(last).Next()
}

// Fail resets the parser to the start an returns the parsed string between the
// start and last. If next it also includes the next in the result.
func (s *state) Fail(start, last *Cursor, next bool) string {
	if last == nil {
		last = start // To prevent nil errors.
	} else if next {
		// Get the mark after the last matching rune.
		last = s.p.Jump(last).Peek()
	}
	s.p.Jump(start) // Reset parser.
	return s.p.Slice(start, last)
}

// End returns a mark to the last successfully parsed rune.
func (s *state) End() *Cursor {
	return s.end
}
