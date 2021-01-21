package recursion

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
)

func Infinite(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{Value: op.Or{And, Value}})
}

func And(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{Infinite, op.MinZero(SP), '+', op.MinZero(SP), Infinite})
}
