package ast_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
	"strconv"
	"testing"
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

func ExampleParser_Expect_parse_node() {
	p, _ := ast.New([]byte("1 <= 2"))
	fmt.Println(p.Expect(func(p *ast.Parser) (*ast.Node, error) {
		digit := ast.Capture{
			Value: parser.CheckRuneFunc(func(r rune) bool {
				return '0' <= r && r <= '9'
			}),
			Convert: func(i string) interface{} {
				v, _ := strconv.Atoi(i)
				return v
			},
		}
		return p.Expect(op.And{digit, parser.CheckString(" <= "), digit})
	}))
	// Output:
	// [-01] [[000] 1, [000] 2] <nil>
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
	// [000] 1 <nil>
	// <nil> <nil>
	// [000] 2 <nil>
}

func ExampleParser_Expect_not() {
	p, _ := ast.New([]byte("bar"))

	_, err := p.Expect(op.Not{Value: "bar"})
	fmt.Println(err)
	_, err = p.Expect(op.Not{
		Value: ast.Capture{
			Value: "baz",
		},
	})
	fmt.Println(err)
	// Output:
	// parse: expected op.Not {bar} but got "bar"
	// <nil>
}

func TestParser_Expect_not(t *testing.T) {
	p, _ := ast.New([]byte("bar\nbaz"))
	any := ast.Capture{
		Value: op.MinZero(op.And{
			op.Not{Value: '\n'},
			parser.CheckRuneFunc(func(r rune) bool {
				return r != parser.EOD
			}),
		}),
	}
	node, err := p.Expect(op.And{
		any,
		'\n',
		any,
	})
	if err != nil {
		t.Error(err)
	}
	if len(node.Children()) != 2 {
		t.Error(node)
	}
}

func ExampleParser_Expect_and() {
	p, _ := ast.New([]byte("1 <= 2"))
	digit := ast.Capture{
		Type: 1,
		Value: parser.CheckRuneFunc(func(r rune) bool {
			return '0' <= r && r <= '9'
		}),
	}

	fmt.Println(p.Expect(op.And{
		digit, parser.CheckString(" <= "), digit,
	}))
	// Output:
	// [-01] [[001] 1, [001] 2] <nil>
}

func ExampleParser_Expect_or() {
	p, _ := ast.New([]byte("data"))

	var (
		d    = ast.Capture{Value: 'd'}
		da   = ast.Capture{Value: "da"}
		data = ast.Capture{Value: "data"}
		a    = ast.Capture{Value: 'a'}
		at   = ast.Capture{Value: "at"}
		ata  = ast.Capture{Value: "ata"}
		t    = ast.Capture{Value: 't'}
	)

	fmt.Println(p.Expect(op.Or{d, da, data}))
	fmt.Println(p.Expect(op.Or{at, a, ata}))
	fmt.Println(p.Expect(op.Or{d, t, op.Not{Value: a}}))
	// Output:
	// [000] d <nil>
	// [000] at <nil>
	// <nil> parse: expected op.Or [{0 100 <nil>} {0 116 <nil>} {{0 97 <nil>}}] but got "a"
}

func TestParser_Expect_and_or(t *testing.T) {
	p, _ := ast.New([]byte("u10FFFF"))
	node, err := p.Expect(
		ast.Capture{
			Value: op.And{
				"u",
				op.Or{
					op.And{
						"10",
						op.Repeat(4,
							op.Or{
								parser.CheckRuneRange('0', '9'),
								parser.CheckRuneRange('A', 'F'),
							},
						),
					},
					op.MinMax(4, 5,
						op.Or{
							parser.CheckRuneRange('0', '9'),
							parser.CheckRuneRange('A', 'F'),
						},
					),
				},
			},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
	if node.ValueString() != "u10FFFF" {
		t.Error(node)
	}
}

func ExampleParser_Expect_xor() {
	p, _ := ast.New([]byte("data"))

	var (
		d    = ast.Capture{Value: 'd'}
		da   = ast.Capture{Value: "da"}
		data = ast.Capture{Value: "data"}
		a    = ast.Capture{Value: 'a'}
		t    = ast.Capture{Value: 't'}
	)

	fmt.Println(p.Expect(op.XOr{d, da, data}))
	fmt.Println(p.Expect(op.XOr{a, t}))
	// Output:
	// <nil> parse: expected op.XOr [{0 100 <nil>} {0 da <nil>} {0 data <nil>}] but got "da"
	// <nil> parse: expected op.XOr [{0 97 <nil>} {0 116 <nil>}] but got "d"
}

func ExampleParser_Expect_range() {
	p, _ := ast.New([]byte("aaa"))
	fmt.Println(p.Expect(ast.Capture{
		Value: op.Min(3, 'a'),
	})) // 3 * 'a'

	p, _ = ast.New([]byte("aaa"))
	fmt.Println(p.Expect(op.Min(4, 'a'))) // err
	// Output:
	// [000] aaa <nil>
	// <nil> parse: expected op.Range {4 -1 97} but got "aaa"
}
