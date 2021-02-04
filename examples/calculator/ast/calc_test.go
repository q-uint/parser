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
	// ["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"]]]]] <nil>
	// ["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"],["MulDiv","*"],["Integer","1"]]]]] <nil>
	// ["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"],["MulDiv","*"],["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"]]]]] <nil>
	// ["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"],["MulDiv","*"],["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"]]]]]]]]] <nil>
	// ["AddSubExpr",[["MulDivExpr",[["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"]]]]],["MulDiv","*"],["AddSubExpr",[["MulDivExpr",[["Integer","1"]]],["AddSub","+"],["MulDivExpr",[["Integer","1"]]]]]]]]] <nil>
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
	// ["Integer","007"] <nil>
	// ["Integer","007"] <nil>
	// ["MulDivExpr",[["Integer","007"]]] <nil>
	// ["AddSubExpr",[["MulDivExpr",[["Integer","007"]]]]] <nil>
}
