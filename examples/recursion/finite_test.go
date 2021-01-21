package recursion

import (
	"fmt"
	"github.com/di-wu/parser/ast"
)

func ExampleFinite() {
	p, _ := ast.New([]byte("0+1 + 10"))
	fmt.Println(p.Expect(ast.Capture{Value: Finite}))
	// Output:
	// [000] [[000] 0, [000] 1, [000] 10] <nil>
}
