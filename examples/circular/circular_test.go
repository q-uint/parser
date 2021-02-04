package circular_test

import (
	"fmt"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/examples/circular"
	"github.com/di-wu/parser/op"
)

func ExampleCircular() {
	p, _ := ast.New([]byte("00001"))
	fmt.Println(p.Expect(ast.Capture{
		TypeStrings: []string{"Number"},
		Value:       circular.Circular,
	}))
	// Output:
	// ["Number","00001"] <nil>
}

func ExampleCircular_function() {
	var circularFunc interface{}
	circularFunc = func(p *ast.Parser) (*ast.Node, error) {
		return p.Expect(op.Or{
			op.And{
				'0',
				circularFunc,
			},
			'1',
		})
	}

	p, _ := ast.New([]byte("00001"))
	fmt.Println(p.Expect(ast.Capture{
		TypeStrings: []string{"Circular"},
		Value:       circularFunc,
	}))
	// Output:
	// ["Circular","00001"] <nil>
}

func ExampleParse() {
	fmt.Println(circular.Parse("00001"))
	fmt.Println(circular.Parse("01001"))
	fmt.Println(circular.Parse("00100"))
	// Output:
	// ["Circular","00001"] <nil>
	// ["Circular","01"] <nil>
	// ["Circular","001"] <nil>
}
