package parser

import (
	"fmt"
	"github.com/di-wu/parser/op"
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
		position: size - 1,
	}
	return &p, nil
}

// Next advances the parser by one rune.
func (p *Parser) Next() *Parser {
	if p.Done() {
		return p
	}

	current, size := utf8.DecodeRune(p.buffer[p.cursor.position+1:])
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

// Peek returns the next cursor without advancing the parser.
func (p *Parser) Peek() *Cursor {
	start := p.Mark()
	defer p.Jump(start)
	return p.Next().Mark()
}

// Jump goes to the position of the given mark.
func (p *Parser) Jump(mark *Cursor) *Parser {
	cursor := *mark
	p.cursor = &cursor
	return p
}

// Slice returns the value in between the two given cursors [start:end]. The end
// value is inclusive!
func (p *Parser) Slice(start *Cursor, end *Cursor) string {
	return string(p.buffer[start.position : end.position+1])
}

// Expect checks whether the buffer contains the given value. It consumes their
// corresponding runes and returns a mark to the last rune of the consumed
// value. It returns an error if can not find a match with the given value.
//
// It currently supports:
//	- rune & string
//	- func(p *Parser) (*Cursor, bool)
//	  (== AnonymousClass)
//	- []interface{}
//	  (== op.And)
//	- operators: op.Not, op.And, op.Or & op.XOr
func (p *Parser) Expect(i interface{}) (*Cursor, error) {
	var end *Cursor // Contains a start that indicates the end (inclusive).
	ok := func(last *Cursor) {
		end = last
		// We jump to the given cursor (last parsed rune) because it is not
		// guaranteed that the already parser did not pass it.
		p.Jump(last).Next()
	}
	fail := func(start, last *Cursor, next bool) string {
		if last == nil {
			last = start // To prevent nil errors.
		} else if next {
			// Get the mark after the last matching rune.
			last = p.Jump(last).Peek()
		}
		p.Jump(start) // Reset parser.
		return p.Slice(start, last)
	}

	// Converting some values for convenience...
	switch v := i.(type) {
	case int:
		i = rune(v)

	case func(p *Parser) (*Cursor, bool):
		i = AnonymousClass(v)
	case Class:
		i = AnonymousClass(v.Check)
	case []interface{}:
		i = op.And(v)
	}

	switch start := p.Mark(); v := i.(type) {
	case rune:
		if p.cursor.Rune != v {
			p.Jump(start)
			return nil, &ExpectedParseError{
				Expected: v, Actual: string(p.cursor.Rune),
			}
		}
		ok(p.Mark())
	case string:
		if v == "" {
			return nil, &ExpectError{
				Message: "can not parse empty string",
			}
		}
		for _, r := range []rune(v) {
			if p.cursor.Rune != r {
				conflict := p.Mark()
				p.Jump(start)
				return nil, &ExpectedParseError{
					Expected: v, Actual: p.Slice(start, conflict),
				}
			}
			ok(p.Mark())
		}

	case AnonymousClass:
		last, passed := v(p)
		if !passed {
			return nil, &ExpectedParseError{
				Expected: v, Actual: fail(start, last, true),
			}
		}
		ok(last)

	case op.Not:
		last, err := p.Expect(v.Value)
		p.Jump(start)
		if err == nil {
			return nil, &ExpectedParseError{
				Expected: v, Actual: p.Slice(start, last),
			}
		}
	case op.And:
		var last *Cursor
		for _, i := range v {
			mark, err := p.Expect(i)
			if err != nil {
				return nil, &ExpectedParseError{
					Expected: v, Actual: fail(start, last, true),
				}
			}
			last = mark
		}
		ok(last)
	case op.Or:
		var last *Cursor
		for _, i := range v {
			mark, err := p.Expect(i)
			if err == nil {
				last = mark
				break
			}
		}
		if last == nil {
			return nil, &ExpectedParseError{
				Expected: v, Actual: fail(start, last, true),
			}
		}
		ok(last)
	case op.XOr:
		var last *Cursor
		for _, i := range v {
			mark, err := p.Expect(i)
			if err == nil {
				if last != nil {
					return nil, &ExpectedParseError{
						Expected: v, Actual: fail(start, mark, false),
					}
				}
				last = mark
				p.Jump(start) // Go back to the start.
			}
		}
		if last == nil {
			return nil, &ExpectedParseError{
				Expected: v, Actual: fail(start, last, true),
			}
		}
		ok(last)
	default:
		return nil, &ExpectError{
			Message: fmt.Sprintf("value of type %T are not supported", v),
		}
	}
	return end, nil
}
