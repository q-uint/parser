package parser_test

import (
	"fmt"
	"github.com/di-wu/parser"
)

func ExampleParser_Current() {
	p, _ := parser.New([]byte("some data"))

	current := p.Current()
	fmt.Printf("%U: %c", current, current)

	// Output:
	// U+0073: s
}

func ExampleParser_Next() {
	p, _ := parser.New([]byte("some data"))

	current := p.Next().Current()
	fmt.Printf("%U: %c", current, current)

	// Output:
	// U+006F: o
}

func ExampleParser_Done() {
	p, _ := parser.New([]byte("_"))

	fmt.Println(p.Next().Done())

	// Output:
	// true
}

func ExampleParser_Expect_rune() {
	p, _ := parser.New([]byte("data"))

	mark, _ := p.Expect('d')
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	_, err := p.Expect('d')
	fmt.Println(err)

	mark, _ = p.Expect('a')
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	mark, _ = p.Expect('t')
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	current := p.Current()
	fmt.Printf("%U: %c\n", current, current)

	fmt.Println(p.Next().Done())

	// Output:
	// U+0064: d
	// parse: expected int32 'd' but got 97
	// U+0061: a
	// U+0074: t
	// U+0061: a
	// true
}

func ExampleParser_Expect_string() {
	p, _ := parser.New([]byte("some data"))

	mark, _ := p.Expect("some")
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	p.Next() // Skip space.

	mark, _ = p.Expect("data")
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	// Output:
	// U+0065: e
	// U+0061: a
}
