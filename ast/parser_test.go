package ast_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"strconv"
)

func ExampleParser_Expect_rune() {
	p, _ := ast.New([]byte("data"))

	fmt.Println(p.Expect('d'))
	fmt.Println(p.Expect('d')) // Invalid.
	fmt.Println(p.Expect('a'))
	fmt.Println(p.Expect('t'))
	fmt.Println(p.Expect('a'))

	// Output:
	// <nil> <nil>
	// <nil> parse: expected int32 100 but got "a"
	// <nil> <nil>
	// <nil> <nil>
	// <nil> <nil>
}

func ExampleParser_Expect_string() {
	p, _ := ast.New([]byte("some data"))

	fmt.Println(p.Expect("some"))
	fmt.Println(p.Expect('_'))
	_, _ = p.Expect(' ') // Skip space.
	fmt.Println(p.Expect("data"))

	// Output:
	// <nil> <nil>
	// <nil> parse: expected int32 95 but got " "
	// <nil> <nil>
}

func ExampleParser_Expect_class() {
	p, _ := ast.New([]byte("1 <= 2"))
	digit := func(p *parser.Parser) (*parser.Cursor, bool) {
		r := p.Current()
		return p.Mark(), '0' <= r && r <= '9'
	}
	lt := func(p *parser.Parser) (*parser.Cursor, bool) {
		var last *parser.Cursor
		for _, r := range []rune("<=") {
			if p.Current() != r {
				return nil, false
			}
			last = p.Mark()
			p.Next()
		}
		return last, true
	}

	fmt.Println(p.Expect(digit))
	_, _ = p.Expect(' ') // Skip space.
	fmt.Println(p.Expect(lt))
	_, _ = p.Expect(' ') // Skip space.
	fmt.Println(p.Expect(digit))

	// Output:
	// <nil> <nil>
	// <nil> <nil>
	// <nil> <nil>
}

func ExampleParser_Expect_capture() {
	p, _ := ast.New([]byte("1 <= 2"))
	digit := ast.Capture{
		Value: parser.CheckRuneFunc(func(r rune) bool {
			return '0' <= r && r <= '9'
		}),
		Convert: func(i string) interface{} {
			v, _ := strconv.Atoi(i)
			return v
		},
	}
	lt := parser.CheckString(" <= ")

	fmt.Println(p.Expect(digit))
	fmt.Println(p.Expect(lt))
	fmt.Println(p.Expect(digit))

	// Output:
	// [000]: 1 <nil>
	// <nil> <nil>
	// [000]: 2 <nil>
}
