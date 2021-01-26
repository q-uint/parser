package elf_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/examples/elf"
)

func ExampleParse_invalid() {
	p, _ := ast.New([]byte{
		// magic, 32-bit, msb, version, padding
		0x7f, 0x45, 0x4c, 0x46, 0x01, 0x02, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// some type
		0xff, 0xff,
	})
	// This will fail since ast uses utf8.DecodeRune to decode runes.
	_, err := elf.Header(p)
	fmt.Printf("0x%x", err.(*parser.ExpectedParseError).String)
	// Output:
	// 0xffff
}

func ExampleParse() {
	header, _ := elf.Parse([]byte{
		// magic, 32-bit, msb, version, padding
		0x7f, 0x45, 0x4c, 0x46, 0x01, 0x02, 0x01, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		// executable file
		0x00, 0x02,
	})
	fmt.Println(header.Type)

	var (
		children = header.Children()

		indent    = children[0]
		iChildren = indent.Children()
		class     = iChildren[0]
		data      = iChildren[0]

		typ = children[1]
	)
	fmt.Println(indent.Type)
	fmt.Printf("\t%x\n", class.Value)
	fmt.Printf("\t%x\n", data.Value)
	fmt.Println(typ.Type)
	fmt.Printf("\t%x\n", typ.Value)
	// Output:
	// 1
	// 2
	// 	01
	// 	01
	// 5
	// 	0002
}
