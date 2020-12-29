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

	mark, _ := alpha(p)
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	// Output:
	// U+0064: d
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

	mark, _ := walrus(p)
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)

	// Output:
	// U+003D: =
}
