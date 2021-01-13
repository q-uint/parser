package im

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
)

// InMemoryParseNode represents a function to parse in-memory nodes.
type InMemoryParseNode func(p *Parser) (*ast.Node, error)

// LoopUp allows for circular references to be used.
type LoopUp struct {
	Key   string
	Table *map[string]interface{}
}

// Parser represents a general purpose in-memory parser.
type Parser struct {
	internal *ast.Parser
	table    map[string]interface{}
}

// New creates a new Parser.
func New(input []byte, table map[string]interface{}) (*Parser, error) {
	internal, err := ast.New(input)
	if err != nil {
		return nil, err
	}
	if table == nil {
		table = make(map[string]interface{})
	}
	return NewFromParser(internal, table)
}

func NewFromParser(ap *ast.Parser, table map[string]interface{}) (*Parser, error) {
	p := Parser{
		internal: ap,
		table:    table,
	}
	p.internal.SetOperator(func(i interface{}) (*ast.Node, error) {
		switch v := i.(type) {
		case LoopUp:
			i, ok := p.table[v.Key]
			if !ok {
				return nil, fmt.Errorf("could not find %s", v.Key)
			}
			return p.Expect(i)

		case InMemoryParseNode:
			return v(&p)

		default:
			return nil, &parser.UnsupportedType{
				Value: i,
			}
		}
	})
	return &p, nil
}

// Expect checks whether the buffer contains the given value.
func (imp *Parser) Expect(i interface{}) (*ast.Node, error) {
	return imp.internal.Expect(ConvertAliases(i))
}

// ConvertAliases converts various default primitive types to aliases for type
// matching.
func ConvertAliases(i interface{}) interface{} {
	switch v := i.(type) {
	case func(p *Parser) (*ast.Node, error):
		return InMemoryParseNode(v)

	default:
		return parser.ConvertAliases(i)
	}
}
