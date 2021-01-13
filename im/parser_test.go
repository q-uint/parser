package im_test

import (
	"fmt"
	"github.com/di-wu/parser"
	"github.com/di-wu/parser/ast"
	"github.com/di-wu/parser/im"
	"github.com/di-wu/parser/op"
)

func Example_loopUp() {
	var table map[string]interface{}
	table = map[string]interface{}{
		"leading0": op.Or{
			op.And{
				'0',
				im.LoopUp{
					Key:   "leading0",
					Table: &table,
				},
			},
			ast.Capture{
				Type:  1,
				Value: parser.CheckRuneRange('1', '9'),
			},
		},
	}
	p, _ := im.New([]byte("000010"), table)
	fmt.Println(p.Expect(ast.Capture{
		Value: table["leading0"],
	}))
	// Output:
	// [000] [[001] 1] <nil>
}
