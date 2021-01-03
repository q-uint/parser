package ast

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/op"
)

// Parser represents a general purpose AST parser.
type Parser struct {
	internal *parser.Parser
}

// New creates a new Parser.
func New(input []byte) (*Parser, error) {
	internal, err := parser.New(input)
	if err != nil {
		return nil, err
	}
	p := Parser{
		internal: internal,
	}
	return &p, err
}

// Expect checks whether the buffer contains the given value.
func (ap *Parser) Expect(i interface{}) (*Node, error) {
	i = ConvertAliases(i)

	p := ap.internal
	switch start := p.Mark(); v := i.(type) {
	case rune, string, parser.AnonymousClass:
		if _, err := p.Expect(v); err != nil {
			return nil, err
		}

	case Capture:
		node, err := ap.Expect(v.Value)
		if node != nil || err != nil {
			return node, err
		}

		if v.Convert != nil {
			return &Node{
				Type:  v.Type,
				Value: v.Convert(p.Slice(start, p.LookBack())),
			}, nil
		}
		return &Node{
			Type:  v.Type,
			Value: p.Slice(start, p.LookBack()),
		}, nil

	case op.Not:
		defer p.Jump(start)
		if _, err := ap.Expect(v.Value); err == nil {
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, p.LookBack()),
			}
		}
	case op.And:
		node := &Node{}
		for _, i := range v {
			n, err := ap.Expect(i)
			if err != nil {
				p.Jump(start)
				return nil, err
			}
			if n != nil {
				node.SetLast(n)
			}
		}
		return node, nil
	case op.Or:
		for _, i := range v {
			node, err := ap.Expect(i)
			if node != nil && err == nil {
				return node, err
			}
			p.Jump(start)
		}
		return nil, &parser.ExpectedParseError{
			Expected: v, Actual: p.Slice(start, p.Mark()),
		}
	case op.XOr:
		var (
			last *parser.Cursor
			node *Node
		)
		for _, i := range v {
			n, err := ap.Expect(i)
			if n != nil && err == nil {
				if node != nil {
					p.Jump(start)
					return nil, &parser.ExpectedParseError{
						Expected: v, Actual: p.Slice(start, last),
					}
				}
				node = n
			}
			last = p.Mark()
			p.Jump(start)
		}
		if node == nil {
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, p.Mark()),
			}
		}
		return node, nil

	case op.Range:
		var (
			count int
			last  *parser.Cursor
			node  = &Node{}
		)
		for {
			n, err := ap.Expect(v.Value)
			if err != nil {
				break
			}
			if n != nil {
				node.SetLast(n)
			}
			last = p.LookBack()
			count++

			if v.Max != -1 && count == v.Max {
				// Break if you have parsed the maximum amount of values.
				// This way count will never be larger than v.Max.
				break
			}
		}
		if count < v.Min {
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, last),
			}
		}

		if node.IsParent() {
			// Only return node if it has children.
			return node, nil
		}

	default:
		return nil, &parser.UnsupportedType{
			Value: i,
		}
	}

	return nil, nil
}

// ConvertAliases converts various default primitive types to aliases for type
// matching.
func ConvertAliases(i interface{}) interface{} {
	switch v := i.(type) {
	case func(p *Parser) (*Node, error):
		return ParseNode(v)

	default:
		return parser.ConvertAliases(i)
	}
}
