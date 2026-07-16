package precedence

import (
	"github.com/q-uint/parser"
	"github.com/q-uint/parser/ast"
	"github.com/q-uint/parser/op"
)

func Value(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		TypeStrings: types,
		Value: op.Or{
			'0',
			op.And{
				parser.CheckRuneRange('1', '9'),
				op.MinZero(parser.CheckRuneRange('0', '9')),
			},
		},
	})
}

const SP = ' '
