package parser_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/is"
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

func TestParser_Slice(t *testing.T) {
	var (
		p, _ = parser.New([]byte("abc"))
		m1   = p.Mark()
		m2   = p.Next().Mark()
		m3   = p.Next().Mark()
		m4   = p.Next().Mark()
	)

	if m1.Rune != 'a' {
		t.Error(m1.Rune)
	}
	if m2.Rune != 'b' {
		t.Error(m2.Rune)
	}
	if m3.Rune != 'c' {
		t.Error(m3.Rune)
	}
	if m4.Rune != parser.EOD {
		t.Error(m4.Rune)
	}

	if s := p.Slice(m1, m3); s != "abc" {
		t.Error(s)
	}
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
	// parse: expected int32 'd' but got "a"
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

func TestParser_Expect_string_err(t *testing.T) {
	p, _ := parser.New([]byte("bar"))

	_, err := p.Expect("baz")
	if err == nil {
		t.Error()
		return
	}

	expected := err.(*parser.ExpectedParseError)
	if expected.Actual != "bar" {
		t.Error(expected.Actual)
	}
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

func testAnonymousClass (p *parser.Parser) (*parser.Cursor, bool) {
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

	mark, err = p.Expect(testAnonymousClass)
	if err != nil {
		t.Error(err)
		return
	}
	if mark.Rune != 'o' {
		t.Error()
	}
}

func TestParser_Expect_class_err(t *testing.T) {
	p, _ := parser.New([]byte("bar"))
	_, err := p.Expect(func(p *parser.Parser) (*parser.Cursor, bool) {
		var last *parser.Cursor
		for _, r := range []rune("baz") {
			if p.Current() != r {
				return last, false
			}
			last = p.Mark()
			p.Next()
		}
		return last, true
	})
	if err == nil {
		t.Error()
		return
	}

	expected := err.(*parser.ExpectedParseError)
	if expected.Actual != "bar" {
		t.Error(expected.Actual)
	}
}

func ExampleParser_Expect_not() {
	p, _ := parser.New([]byte("bar"))

	_, err := p.Expect(is.Not{Value: "baz"})
	fmt.Println(err)
	_, err = p.Expect(is.Not{Value: "bar"})
	fmt.Println(err)

	// Output:
	// <nil>
	// parse: expected is.Not {"bar"} but got "bar"
}
