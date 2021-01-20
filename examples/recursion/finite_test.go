package recursion

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFinite() {
	p, _ := ast.New([]byte("x&&x&&x"))
	fmt.Println(p.Expect(ast.Capture{Value: Finite}))
	// Output:
	// [000] [[000] x, [000] x, [000] x] <nil>
}
