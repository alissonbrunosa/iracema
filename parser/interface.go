package parser

import (
	"io"
	"iracema/ast"
)

func Parse(input io.Reader) (file *ast.File, err error) {
	p := new(parser)
	p.init(input)

	file = p.parse()
	return file, p.errors.Err()
}
