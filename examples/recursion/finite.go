package recursion

import (
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
)

func Finite(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{Value, op.MinZero(op.And{op.MinZero(SP), '+', op.MinZero(SP), Finite})})
}
