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
		Value: circular.Circular,
	}))
	// Output:
	// [000] 00001 <nil>
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
		Value: circularFunc,
	}))
	// Output:
	// [000] 00001 <nil>
}

func ExampleParse() {
	fmt.Println(circular.Parse("00001"))
	fmt.Println(circular.Parse("01001"))
	fmt.Println(circular.Parse("00100"))
	// Output:
	// [000] 00001 <nil>
	// [000] 01 <nil>
	// [000] 001 <nil>
}
