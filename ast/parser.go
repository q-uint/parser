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
		// Just check if it matches.
		if _, err := p.Expect(v); err != nil {
			return nil, err
		}

	case ParseNode:
		node, err := v(ap)
		if err != nil {
			p.Jump(start)
			return nil, err
		}
		return node, nil

	case Capture:
		node, err := ap.Expect(v.Value)
		if err != nil {
			p.Jump(start)
			return nil, err
		}
		if node != nil {
			// Return the node.
			if node.Type == -1 {
				node.Type = v.Type
			}
			return node, nil
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
			// Return error if match is found.
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, p.LookBack()),
			}
		}
	case op.And:
		node := &Node{Type: -1}
		for _, i := range v {
			n, err := ap.Expect(i)
			if err != nil {
				p.Jump(start)
				return nil, err
			}
			if n != nil {
				if n.Type == -1 {
					node.Adopt(n)
				} else {
					node.SetLast(n)
				}
			}
		}
		return node, nil
	case op.Or:
		// To keep track whether we encountered a valid value, node or not.
		var hit bool
		for _, i := range v {
			node, err := ap.Expect(i)
			if err == nil {
				hit = true
				if node != nil {
					// Return node if found.
					return node, nil
				}
				break
			}
			p.Jump(start)
		}
		if !hit {
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, p.Mark()),
			}
		}
	case op.XOr:
		var (
			node *Node
			last *parser.Cursor
		)
		for _, i := range v {
			n, err := ap.Expect(i)
			if err != nil {
				p.Jump(start)
				continue
			}
			if last != nil {
				// We already got a match.
				p.Jump(start)
				return nil, &parser.ExpectedParseError{
					Expected: v, Actual: p.Slice(start, last),
				}
			}
			last = p.Mark()
			node = n
		}
		if last == nil {
			p.Jump(start)
			return nil, &parser.ExpectedParseError{
				Expected: v, Actual: p.Slice(start, p.Mark()),
			}
		}
		if node != nil {
			return node, nil
		}

	case op.Range:
		var (
			count int
			last  *parser.Cursor
			node  = &Node{Type: -1}
		)
		for {
			n, err := ap.Expect(v.Value)
			if err != nil {
				break
			}
			if n != nil {
				if n.Type == -1 {
					node.Adopt(n)
				} else {
					node.SetLast(n)
				}
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
			p.Jump(start)
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