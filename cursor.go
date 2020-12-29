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
