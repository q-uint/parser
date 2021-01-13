package circular

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/im"
	"github.com/di-wu/parser/op"
)

func Circular(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		op.And{
			'0',
			Circular,
		},
		'1',
	})
}

var table map[string]interface{}

func Parse(input string) (*ast.Node, error) {
	p, err := NewCircularParser([]byte(input))
	if err != nil {
		return nil, err
	}
	return p.Expect(ast.Capture{
		Value: table["circular"],
	})
}

func NewCircularParser(input []byte) (*im.Parser, error) {
	table = map[string]interface{}{
		"circular": op.Or{
			op.And{
				'0',
				im.LoopUp{
					Key:   "circular",
					Table: &table,
				},
			},
			'1',
		},
	}
	p, err := im.New(input, table)
	if err != nil {
		return nil, err
	}
	return p, nil
}
