package elf

import (
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/op"
)

func Parse(input []byte) (*ast.Node, error) {
	p, err := NewELFParser(input)
	if err != nil {
		return nil, err
	}
	return Header(p)
}

func NewELFParser(input []byte) (*ast.Parser, error) {
	p, err := ast.New(input)
	if err != nil {
		return nil, err
	}
	p.DecodeRune(func(p []byte) (rune, int) {
		if len(p) == 0 {
			return rune(-1), 0
		}
		return rune(p[0]), 1
	})
	return p, nil
}

func Header(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: 1,
		Value: op.And{
			Indent, Type,
			// etc.
		},
	})
}

func Indent(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type: 2,
		Value: op.And{
			MagicNumber,
			Class,
			Data,
			Version,
			op.Repeat(9, 0x00),
		},
	})
}

func MagicNumber(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(op.And{
		0x7f, 0x45, 0x4c, 0x46,
	})
}

func Class(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:  3,
		Value: parser.CheckRuneRange(0x00, 0x02),
	})
}

func Data(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:  4,
		Value: parser.CheckRuneRange(0x00, 0x02),
	})
}
func Version(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(0x01)
}

func Type(p *ast.Parser) (*ast.Node, error) {
	return p.Expect(ast.Capture{
		Type:  5,
		Value: op.Repeat(2, parser.CheckRuneRange(0x00, 0xff)),
	})
}
