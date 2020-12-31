package calc_test

import (
	"fmt"
	calc "github.com/di-wu/parser/examples/calculator"
)

func ExampleParse() {
	fmt.Println(calc.Parse("1 + 1"))
	fmt.Println(calc.Parse("1 + 1 * 1"))
	fmt.Println(calc.Parse("1 + 1 * 1 + 1"))
	fmt.Println(calc.Parse("1 + 1 * (1 + 1)"))
	fmt.Println(calc.Parse("(1 + 1) * (1 + 1)"))

	// Output:
	// (1 + 1) <nil>
	// (1 + (1 * 1)) <nil>
	// (1 + (1 * 1) + 1) <nil>
	// (1 + (1 * (1 + 1))) <nil>
	// ((1 + 1) * (1 + 1)) <nil>
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

func ExampleCalculatorParser_Integer() {
	p, _ := calc.New("007")
	fmt.Println(p.Integer())

	// Output:
	// 7 <nil>
}
