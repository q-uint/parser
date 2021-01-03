package parser

import (
	"github.com/di-wu/parser/op"
	"unicode/utf8"
)

// EOD indicates the End Of (the) Data.
const EOD = 1<<31 - 1

// Parser represents a general purpose parser.
type Parser struct {
	buffer []byte
	cursor *Cursor

	converter func(i interface{}) interface{}
	operator  func(i interface{}) (*Cursor, error)
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

// SetConverter allows you to add additional (prioritized) converters to the
// parser. e.g. convert aliases to other types or overwrite defaults.
func (p *Parser) SetConverter(c func(i interface{}) interface{}) {
	p.converter = c
}

// SetOperator allows you to support additional (prioritized) operators.
// Should return an UnsupportedType error if the given value is not supported.
func (p *Parser) SetOperator(o func(i interface{}) (*Cursor, error)) {
	p.operator = o
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

// LookBack returns the previous cursor without decreasing the parser.
func (p *Parser) LookBack() *Cursor {
	if p.cursor.position == 0 || p.Done() {
		// Not possible to go back
		return p.Mark()
	}

	// We don't know the size of the previous rune... 1 or more?
	previous, size := utf8.DecodeRune(p.buffer[p.cursor.position-1:])
	for i := 2; previous == utf8.RuneError; i++ {
		previous, size = utf8.DecodeRune(p.buffer[p.cursor.position-i:])
	}
	return &Cursor{
		Rune:     previous,
		position: p.cursor.position - size,
	}
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
	if end == nil { // Just to be sure...
		end = start
	}
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
	state := state{p: p}

	i = ConvertAliases(i)
	if p.converter != nil {
		// Can undo previous conversions!
		i = p.converter(i)
	}

	if p.operator != nil {
		// Takes priority over default values. If an unsupported error is
		// returned we can check if one of the predefined types match.
		mark, err := p.operator(i)
		if _, ok := err.(*UnsupportedType); !ok {
			return mark, err
		}
	}
	switch start := p.Mark(); v := i.(type) {
	case rune:
		if p.cursor.Rune != v {
			p.Jump(start)
			return nil, &ExpectedParseError{
				Expected: v, Actual: string(p.cursor.Rune),
			}
		}
		state.Ok(p.Mark())
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
			state.Ok(p.Mark())
		}

	case AnonymousClass:
		last, passed := v(p)
		if !passed {
			return nil, &ExpectedParseError{
				Expected: v, Actual: state.Fail(start, last, true),
			}
		}
		state.Ok(last)

	case op.Not:
		defer p.Jump(start)
		if last, err := p.Expect(v.Value); err == nil {
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
					Expected: v, Actual: state.Fail(start, last, true),
				}
			}
			last = mark
		}
		state.Ok(last)
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
				Expected: v, Actual: state.Fail(start, last, true),
			}
		}
		state.Ok(last)
	case op.XOr:
		var last *Cursor
		for _, i := range v {
			mark, err := p.Expect(i)
			if err == nil {
				if last != nil {
					return nil, &ExpectedParseError{
						Expected: v, Actual: state.Fail(start, mark, false),
					}
				}
				last = mark
				p.Jump(start) // Go back to the start.
			}
		}
		if last == nil {
			return nil, &ExpectedParseError{
				Expected: v, Actual: state.Fail(start, last, true),
			}
		}
		state.Ok(last)

	case op.Range:
		var (
			count int
			last  *Cursor
		)
		for {
			mark, err := p.Expect(v.Value)
			if err != nil {
				break
			}
			last = mark
			count++

			if v.Max != -1 && count == v.Max {
				// Break if you have parsed the maximum amount of values.
				// This way count will never be larger than v.Max.
				break
			}
		}
		if count < v.Min {
			return nil, &ExpectedParseError{
				Expected: v, Actual: state.Fail(start, last, true),
			}
		}
		state.Ok(last)

	default:
		return nil, &UnsupportedType{
			Value: i,
		}
	}
	return state.End(), nil
}

// ConvertAliases converts various default primitive types to aliases for type
// matching.
func ConvertAliases(i interface{}) interface{} {
	switch v := i.(type) {
	case int:
		return rune(v)

	case func(p *Parser) (*Cursor, bool):
		return AnonymousClass(v)
	case Class:
		return AnonymousClass(v.Check)

	case []interface{}:
		return op.And(v)

	default:
		return i
	}
}
