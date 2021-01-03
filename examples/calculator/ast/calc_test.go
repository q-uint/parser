package calc_ast_test

import (
	"fmt"
	"github.com/di-wu/parser/ast"
	calc "github.com/di-wu/parser/examples/calculator/ast"
)

func ExampleParse() {
	fmt.Println(calc.Parse("1 + 1"))
	fmt.Println(calc.Parse("1 + 1 * 1"))
	fmt.Println(calc.Parse("1 + 1 * 1 + 1"))
	fmt.Println(calc.Parse("1 + 1 * (1 + 1)"))
	fmt.Println(calc.Parse("(1 + 1) * (1 + 1)"))

	// Output:
	// [001] [[002] [[005] 1], [003] 43, [002] [[005] 1]] <nil>
	// [001] [[002] [[005] 1], [003] 43, [002] [[005] 1, [004] 42, [005] 1]] <nil>
	// [001] [[002] [[005] 1], [003] 43, [002] [[005] 1, [004] 42, [005] 1], [003] 43, [002] [[005] 1]] <nil>
	// [001] [[002] [[005] 1], [003] 43, [002] [[005] 1, [004] 42, [001] [[002] [[005] 1], [003] 43, [002] [[005] 1]]]] <nil>
	// [001] [[002] [[001] [[002] [[005] 1], [003] 43, [002] [[005] 1]], [004] 42, [001] [[002] [[005] 1], [003] 43, [002] [[005] 1]]]] <nil>
}

func ExampleEvaluate() {
	fmt.Println(calc.Evaluate("1 + 1"))
	fmt.Println(calc.Evaluate("1 + 1 * 1"))
	fmt.Println(calc.Evaluate("1 + 1 * 1 + 1"))
	fmt.Println(calc.Evaluate("1 + 1 * (1 + 1)"))
	fmt.Println(calc.Evaluate("(1 + 1) * (1 + 1)"))

	// Output:
	// 2 <nil>
	// 2 <nil>
	// 3 <nil>
	// 3 <nil>
	// 4 <nil>
}


func ExampleInteger() {
	p := func(s string) *ast.Parser {
		p, _ := ast.New([]byte(s))
		return p
	}
	fmt.Println(p("007").Expect(calc.Integer))
	fmt.Println(p("007").Expect(calc.Factor))
	fmt.Println(p("007").Expect(calc.MulDiv))
	fmt.Println(p("007").Expect(calc.AddSub))

	// Output:
	// [005] 7 <nil>
	// [005] 7 <nil>
	// [002] [[005] 7] <nil>
	// [001] [[002] [[005] 7]] <nil>
}
