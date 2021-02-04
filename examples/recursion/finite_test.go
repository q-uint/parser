package recursion

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFinite() {
	p, _ := ast.New([]byte("0+1 + 10"))
	fmt.Println(p.Expect(ast.Capture{
		TypeStrings: []string{"Finite"},
		Value:       Finite,
	}))
	// Output:
	// ["Finite",[["Value","0"],["Value","1"],["Value","10"]]] <nil>
}
