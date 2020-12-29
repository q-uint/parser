package parser_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"testing"
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

func ExampleParser_Expect_class() {
	p, _ := parser.New([]byte("1 <= 2"))
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

	mark, _ := p.Expect(digit)
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	p.Next() // Skip space.

	mark, _ = p.Expect(lt)
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	p.Next() // Skip space.

	mark, _ = p.Expect(digit)
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	// Output:
	// U+0031: 1
	// U+003D: =
	// U+0032: 2
}

type testClass struct{}

func (t testClass) Check(p *parser.Parser) (*parser.Cursor, bool) {
	r := p.Current()
	return p.Mark(), 'a' <= r && r <= 'z'
}

func TestParser_Expect_class(t *testing.T) {
	p, _ := parser.New([]byte("some data"))
	mark, err := p.Expect(testClass{})
	if err != nil {
		t.Error(err)
		return
	}
	if mark.Rune != 's' {
		t.Error()
	}
}
