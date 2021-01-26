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
	// parse conflict [00:002]: expected op.Not !"bar" but got "bar"
	// <nil>
}

func ExampleEnsure() {
	p, _ := parser.New([]byte("bar"))

	fmt.Println(p.Expect(op.Ensure{Value: "ba"}))
	fmt.Println(p.Expect(op.Ensure{Value: "baz"}))
	fmt.Println(p.Expect("bar"))
	// Output:
	// <nil> <nil>
	// <nil> parse conflict [00:002]: expected string "baz" but got "bar"
	// U+0072: r <nil>
}

func ExampleAnd() {
	p, _ := parser.New([]byte("foo bar baz"))

	_, err := p.Expect(op.And{"foo", ' ', "bar", '_'})
	fmt.Println(err)
	_, err = p.Expect(op.And{"foo", ' ', "bar", ' ', "baz"})
	fmt.Println(err)
	// Output:
	// parse conflict [00:007]: expected op.And and["foo" ' ' "bar" '_'] but got "foo bar "
	// <nil>
}

func ExampleOr() {
	p, _ := parser.New([]byte("data"))

	mark, _ := p.Expect(op.Or{'d', "da", "data"})
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	mark, _ = p.Expect(op.Or{"at", 'a', "ata"})
	fmt.Printf("%U: %c\n", mark.Rune, mark.Rune)
	fmt.Println(p.Expect(op.Or{'d', 't', op.Not{'a'}}))
	// Output:
	// U+0064: d
	// U+0074: t
	// <nil> parse conflict [00:003]: expected op.Or or['d' 't' !'a'] but got 'a'
}

func ExampleXOr() {
	p, _ := parser.New([]byte("data"))

	_, err := p.Expect(op.XOr{'d', "da", "data"})
	fmt.Println(err)
	_, err = p.Expect(op.XOr{'a', 't'})
	fmt.Println(err)
	// Output:
	// parse conflict [00:001]: expected op.XOr xor['d' "da" "data"] but got "da"
	// parse conflict [00:000]: expected op.XOr xor['a' 't'] but got 'd'
}
