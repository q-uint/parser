package parser

// Cursor allows you to record your current position so you can return to it
// later. Keeps track of its own position in the buffer of the parser.
type Cursor struct {
	// Rune is the value that the cursor points at.
	Rune rune
	// The position in the buffer.
	position int
}
