package op_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/op"
)

func ExampleMin() {
	p, _ := parser.New([]byte("aaa"))
	start := p.Mark() // Mark to first 'a'.

	fmt.Println(p.Expect(op.Min(3, 'a'))) // 3 * 'a'

	// Reset parser to start.
	p.Jump(start)
	fmt.Println(p.Expect(op.Min(4, 'a'))) // err
	// Output:
	// U+0061 a <nil>
	// <nil> parse: expected op.Range {4 -1 97} but got "aaa"
}

func ExampleMinZero() {
	p, _ := parser.New([]byte("aaab"))

	fmt.Println(p.Expect(op.MinZero('a'))) // 3 * 'a'
	fmt.Println(p.Expect(op.MinZero('b'))) // 1 * 'b'
	fmt.Println(p.Expect(op.MinZero('c'))) // 0 * 'c'
	// Output:
	// U+0061 a <nil>
	// U+0062 b <nil>
	// <nil> <nil>
}

func ExampleMinOne() {
	p, _ := parser.New([]byte("aaab"))

	fmt.Println(p.Expect(op.MinOne('a'))) // 3 * 'a'
	fmt.Println(p.Expect(op.MinOne('b'))) // 1 * 'b'
	fmt.Println(p.Expect(op.MinOne('c'))) // err
	// Output:
	// U+0061 a <nil>
	// U+0062 b <nil>
	// <nil> parse: expected op.Range {1 -1 99} but got "b"
}

func ExampleMinMax() {
	p, _ := parser.New([]byte("aaab"))

	fmt.Println(p.Expect(op.MinMax(1, 2, 'a')))  // 2 * 'a'
	fmt.Println(p.Expect(op.MinMax(1, 3, 'a')))  // 1 * 'a'
	fmt.Println(p.Expect(op.MinMax(0, 1, 'b')))  // 1 * 'b'
	fmt.Println(p.Expect(op.MinMax(1, -1, 'c'))) // err
	// Output:
	// U+0061 a <nil>
	// U+0061 a <nil>
	// U+0062 b <nil>
	// <nil> parse: expected op.Range {1 -1 99} but got "b"
}

func ExampleOptional() {
	p, _ := parser.New([]byte("ac"))

	fmt.Println(p.Expect(op.Optional('a'))) // 1 * 'a'
	fmt.Println(p.Expect(op.Optional('b'))) // 0 * 'a'
	fmt.Println(p.Expect(op.Optional('c'))) // 1 * 'a'
	// Output:
	// U+0061 a <nil>
	// <nil> <nil>
	// U+0063 c <nil>
}

func ExampleRepeat() {
	p, _ := parser.New([]byte("abbbc"))

	fmt.Println(p.Expect(op.Repeat(1, 'a'))) // 1 * 'a'
	fmt.Println(p.Expect(op.Repeat(2, 'b'))) // 2 * b'
	fmt.Println(p.Expect(op.Repeat(1, 'c'))) // err
	// Output:
	// U+0061 a <nil>
	// U+0062 b <nil>
	// <nil> parse: expected op.Range {1 1 99} but got "b"
}
