package ast

import (
	"github.com/di-wu/parser"
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
	default:
		_ = start
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
