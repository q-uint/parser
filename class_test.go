package parser_test

import (
	"fmt"
	"github.com/di-wu/parser"
)

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
