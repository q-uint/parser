package calc

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/op"
	"strconv"
)

func Parse(s string) (*Node, error) {
	p, err := New(s)
	if err != nil {
		return nil, err
	}
	return p.AddSub()
}

func Evaluate(s string) (int, error) {
	n, err := Parse(s)
	if err != nil {
		return 0, err
	}
	return EvaluateNode(n), nil
}

func EvaluateNode(n *Node) int {
	if len(n.Children) == 0 {
		return n.Value
	}
	value := EvaluateNode(n.Children[0])
	for i := range n.Operators {
		switch n.Operators[i] {
		case '+':
			value += EvaluateNode(n.Children[i+1])
		case '-':
			value -= EvaluateNode(n.Children[i+1])
		case '*':
			value *= EvaluateNode(n.Children[i+1])
		case '/':
			value /= EvaluateNode(n.Children[i+1])
		}
	}
	return value
}

type CalculatorParser struct {
	p *parser.Parser
}

// Node is a simple node representation.
//
// Value Node:
//	1 = {1}
// Operator Node:
//	1 * (1 + 2) = {['*'] [{1}, {['+', [1, 2]]}]}
type Node struct {
	Value int // If leaf node.

	Operators []rune  // An list of operators like '+', '-', '*' or '/'.
	Children  []*Node // The len(Children) should be len(Operators)+1
}

func (n *Node) String() string {
	if len(n.Children) == 0 {
		return fmt.Sprintf("%d", n.Value)
	}
	s := fmt.Sprintf("(%s", n.Children[0])
	for i := range n.Operators {
		s += fmt.Sprintf(" %s %s", string(n.Operators[i]), n.Children[i+1])
	}
	return fmt.Sprintf("%s)", s)
}

func New(s string) (*CalculatorParser, error) {
	p, err := parser.New([]byte(s))
	if err != nil {
		return nil, err
	}
	return &CalculatorParser{p: p}, nil
}

func (cp *CalculatorParser) AddSub() (*Node, error) {
	return cp.operator(cp.MulDiv, op.Or{'+', '-'})
}

func (cp *CalculatorParser) MulDiv() (*Node, error) {
	return cp.operator(cp.Factor, op.Or{'*', '/'})
}

func (cp *CalculatorParser) operator(f func() (*Node, error), or op.Or) (*Node, error) {
	// f
	left, err := f()
	if err != nil {
		return nil, err
	}

	var (
		operators = make([]rune, 0)
		children  = []*Node{
			left,
		}
	)

	// ((or) f)*
	for start := cp.p.Mark(); ; {
		// (or) f
		cp.Space() // Consume space.
		sign, err := cp.p.Expect(or)
		if err != nil {
			break
		}
		cp.Space() // Consume space.
		right, err := f()
		if err != nil {
			// Don't forget to reset!
			cp.p.Jump(start)
			break
		}
		// Only add to the slices if both parsed correctly.
		operators = append(operators, sign.Rune)
		children = append(children, right)
	}

	if len(children) == 1 {
		return children[0], nil
	}

	return &Node{
		Operators: operators,
		Children:  children,
	}, nil
}

func (cp *CalculatorParser) Factor() (*Node, error) {
	// Integer
	i, err := cp.Integer()
	if err == nil {
		return &Node{Value: i}, nil
	}

	// LPAREN AddSub RPAREN
	if _, err := cp.p.Expect('('); err != nil {
		return nil, err
	}
	cp.Space() // Consume space.
	addSub, err := cp.AddSub()
	if err != nil {
		return nil, err
	}
	cp.Space() // Consume space.
	if _, err := cp.p.Expect(')'); err != nil {
		return nil, err
	}
	return addSub, err
}

func (cp *CalculatorParser) Integer() (int, error) {
	start := cp.p.Mark()
	last, err := cp.p.Expect(op.MinOne(
		parser.CheckRuneFunc(func(r rune) bool {
			return '0' <= r && r <= '9'
		}),
	))
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(cp.p.Slice(start, last))
}

// Space consumes all the spaces.
func (cp *CalculatorParser) Space() {
	_, _ = cp.p.Expect(op.MinZero(' '))
}
