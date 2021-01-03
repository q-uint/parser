package op_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/op"
)

func ExampleNot() {
	p, _ := parser.New([]byte("bar"))

	_, err := p.Expect(op.Not{Value: "bar"})
	fmt.Println(err)
	_, err = p.Expect(op.Not{Value: "baz"})
	fmt.Println(err)

	// Output:
	// parse: expected op.Not {bar} but got "bar"
	// <nil>
}

func ExampleAnd() {
	p, _ := parser.New([]byte("foo bar baz"))

	_, err := p.Expect(op.And{"foo", ' ', "bar", '_'})
	fmt.Println(err)
	_, err = p.Expect(op.And{"foo", ' ', "bar", ' ', "baz"})
	fmt.Println(err)

	// Output:
	// parse: expected op.And [foo 32 bar 95] but got "foo bar "
	// <nil>
}

func ExampleOr() {
	p, _ := parser.New([]byte("data"))

	mark, _ := p.Expect(op.Or{'d', "da", "data"})
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	mark, _ = p.Expect(op.Or{"at", 'a', "ata"})
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	_, err := p.Expect(op.Or{'d', 't', op.Not{'a'}})
	fmt.Println(err)

	// Output:
	// U+0064: d
	// U+0074: t
	// parse: expected op.Or [100 116 {97}] but got "a"
}

func ExampleXOr() {
	p, _ := parser.New([]byte("data"))

	_, err := p.Expect(op.XOr{'d', "da", "data"})
	fmt.Println(err)
	_, err = p.Expect(op.XOr{'a', 't'})
	fmt.Println(err)

	// Output:
	// parse: expected op.XOr [100 da data] but got "da"
	// parse: expected op.XOr [97 116] but got "d"
}
