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
	// [001] [[002] [[000] 5]] <nil>
	// [001] [[002] [[000] 5], [002] [[000] 2]] <nil>
	// [001] [[002] [[000] 5], [002] [[000] 2, [000] 0]] <nil>
}
