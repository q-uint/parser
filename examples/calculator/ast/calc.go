package calc_ast

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"strconv"
)

func Parse(s string) (*ast.Node, error) {
	p, err := ast.New([]byte(s))
	if err != nil {
		return nil, err
	}
	return p.Expect(AddSub)
}

func Evaluate(s string) (int, error) {
	n, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return EvaluateNode(n), nil
}

func EvaluateNode(n *ast.Node) int {
	var value int
	switch n.Type {
	case 1, 2:
		var sign rune
		for _, c := range n.Children() {
			if c.Type == 3 || c.Type == 4 {
				sign = c.Value.(rune)
			} else {
				switch sign {
				case '+':
					value += EvaluateNode(c)
				case '-':
					value -= EvaluateNode(c)
				case '*':
					value *= EvaluateNode(c)
				case '/':
					value /= EvaluateNode(c)
				default:
					value = EvaluateNode(c)
				}
			}
		}
	case 5:
		return n.Value.(int)
	}
	return value
}

func AddSub(p *ast.Parser) (*ast.Node, error) {
	return operator(p,
		1, MulDiv,
		ast.Capture{
			Type:        3,
			TypeStrings: types,
			Value:       op.Or{'+', '-'},
			Convert: func(i string) interface{} {
				return rune(i[0])
			},
		},
	)
}

func MulDiv(p *ast.Parser) (*ast.Node, error) {
	return operator(p,
		2, Factor,
		ast.Capture{
			Type:        4,
			TypeStrings: types,
			Value:       op.Or{'*', '/'},
			Convert: func(i string) interface{} {
				return rune(i[0])
			},
		},
	)
}

func operator(p *ast.Parser, typ int, f func(p *ast.Parser) (*ast.Node, error), or ast.Capture) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        typ,
		TypeStrings: types,
		Value: op.And{
			f,
			op.MinZero(op.And{
				Space, or, Space, f,
			}),
		},
	})
}

func Factor(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.Or{
		Integer,
		op.And{'(', Space, AddSub, Space, ')'},
	})
}

func Integer(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:        5,
		TypeStrings: types,
		Value: op.MinOne(
			parser.CheckRuneFunc(func(r rune) bool {
				return '0' <= r && r <= '9'
			}),
		),
		Convert: func(i string) interface{} {
			v, _ := strconv.Atoi(i)
			return v
		},
	})
}

// Space consumes all the spaces.
func Space(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.MinZero(' '))
}

var types = []string{
	"UNKNOWN",

	"AddSubExpr",
	"MulDivExpr",
	"AddSub",
	"MulDiv",
	"Integer",
}
