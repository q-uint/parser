package precedence

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func Example() {
	p := func(s string) *ast.Parser {
		p, _ := ast.New([]byte(s))
		return p
	}
	fmt.Println(Plus(p("5")))
	fmt.Println(Plus(p("5 + 2")))
	fmt.Println(Plus(p("5 + 2 * 0")))
	// Output:
	// ["Plus",[["Mult",[["Value","5"]]]]] <nil>
	// ["Plus",[["Mult",[["Value","5"]]],["Mult",[["Value","2"]]]]] <nil>
	// ["Plus",[["Mult",[["Value","5"]]],["Mult",[["Value","2"],["Value","0"]]]]] <nil>
}
