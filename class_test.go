package parser_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"strconv"
	"testing"
)

func ExampleCheckRuneCI() {
	p, _ := parser.New([]byte("Ee"))
	fmt.Println(p.Expect(parser.CheckRune('E')))
	fmt.Println(p.Expect(parser.CheckRuneCI('e')))
	// Output:
	// U+0045: E <nil>
	// U+0065: e <nil>
}

func ExampleCheckStringCI() {
	p, _ := parser.New([]byte("Ee"))
	fmt.Println(p.Expect(parser.CheckStringCI("ee")))
	// Output:
	// U+0065: e <nil>
}

func ExampleCheckInteger() {
	p, _ := parser.New([]byte("-0001 something else"))
	fmt.Println(p.Check(parser.CheckInteger(-1, false)))
	fmt.Println(p.Check(parser.CheckInteger(-1, true)))
	// Output:
	// <nil> false
	// U+0031: 1 true
}

func ExampleCheckIntegerRange() {
	p, _ := parser.New([]byte("12445"))
	fmt.Println(p.Check(parser.CheckIntegerRange(12, 12345, false)))
	fmt.Println(p.Check(parser.CheckIntegerRange(10000, 54321, false)))

	p0, _ := parser.New([]byte("00012"))
	fmt.Println(p0.Check(parser.CheckIntegerRange(12, 12345, false)))
	fmt.Println(p0.Check(parser.CheckIntegerRange(12, 12345, true)))
	// Output:
	// <nil> false
	// U+0035: 5 true
	// <nil> false
	// U+0032: 2 true
}

func TestCheckIntegerRange(t *testing.T) {
	p := func(i int) *parser.Parser {
		p, _ := parser.New([]byte(strconv.Itoa(i)))
		return p
	}

	for i := 0; i < 12; i++ {
		if v, ok := p(i).Check(parser.CheckIntegerRange(12, 12345, false)); ok {
			t.Error(i, v)
		}
	}
	for i := 12; i <= 12345; i++ {
		p := p(i)
		if _, err := p.Expect(parser.CheckIntegerRange(12, 12345, false)); err != nil {
			t.Error(i, err)
		}
		if _, err := p.Expect(parser.EOD); err != nil {
			t.Error(i, err)
		}
	}
}

func ExampleAnonymousClass_rune() {
	p, _ := parser.New([]byte("data"))
	alpha := func(p *parser.Parser) (*parser.Cursor, bool) {
		r := p.Current()
		return p.Mark(),
			'A' <= r && r <= 'Z' ||
				'a' <= r && r <= 'z'
	}

	fmt.Println(alpha(p))
	// Output:
	// U+0064: d true
}

func ExampleAnonymousClass_string() {
	p, _ := parser.New([]byte(":="))
	walrus := func(p *parser.Parser) (*parser.Cursor, bool) {
		var last *parser.Cursor
		for _, r := range []rune(":=") {
			if p.Current() != r {
				return nil, false
			}
			last = p.Mark()
			p.Next()
		}
		return last, true
	}

	fmt.Println(walrus(p))
	// Output:
	// U+003D: = true
}

func ExampleAnonymousClass_error() {
	p, _ := parser.New([]byte("0"))
	lower := func(p *parser.Parser) (*parser.Cursor, bool) {
		r := p.Current()
		return p.Mark(), 'a' <= r && r <= 'z'
	}

	fmt.Println(lower(p))
	// Output:
	// U+0030: 0 false
}
