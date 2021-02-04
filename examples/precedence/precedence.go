package precedence

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
)

func Plus(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        1,
		TypeStrings: types,
		Value: op.And{
			Mult,
			op.MinZero(op.And{
				op.MinZero(SP),
				'+',
				op.MinZero(SP),
				Mult,
			}),
		},
	})
}

func Mult(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        2,
		TypeStrings: types,
		Value: op.And{
			Rule,
			op.MinZero(op.And{
				op.MinZero(SP),
				'*',
				op.MinZero(SP),
				Rule,
			}),
		},
	})
}

func Rule(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		Value,
		op.And{
			'(',
			op.MinZero(SP),
			Plus,
			op.MinZero(SP),
			')',
		},
	})
}

var types = []string{
	"Value",
	"Plus",
	"Mult",
}
