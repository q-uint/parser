package parser

import (
	"unicode/utf8"
)

// EOD indicates the End Of (the) Data.
const EOD = 1<<31 - 1

// Parser represents a general purpose parser.
type Parser struct {
	buffer []byte
	cursor *Cursor
}

// New creates a new Parser.
func New(input []byte) (*Parser, error) {
	p := Parser{
		buffer: input,
	}

	current, size := utf8.DecodeRune(p.buffer)
	if size == 0 {
		// Nothing got decoded.
		return nil, &InitError{
			Message: "failed to scan the first rune",
		}
	}

	p.cursor = &Cursor{
		Rune:     current,
		position: size,
	}
	return &p, nil
}

// Next advances the parser by one rune.
func (p *Parser) Next() *Parser {
	if p.Done() {
		return p
	}

	current, size := utf8.DecodeRune(p.buffer[p.cursor.position:])
	if size == 0 {
		// Nothing got decoded.
		current = EOD
	}

	p.cursor.Rune = current
	p.cursor.position += size

	return p
}

// Current returns the value to which the cursor is pointing at.
func (p *Parser) Current() rune {
	return p.cursor.Rune
}

// Done checks whether the parser is done parsing.
func (p *Parser) Done() bool {
	return p.cursor.Rune == EOD
}

// Mark returns a copy of the current cursor.
func (p *Parser) Mark() *Cursor {
	mark := *p.cursor
	return &mark
}

// Jump goes to the position of the given mark.
func (p *Parser) Jump(mark *Cursor) *Parser {
	cursor := *mark
	p.cursor = &cursor
	return p
}
