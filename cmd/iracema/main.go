package main

import (
	"flag"
	"fmt"
	"iracema/ast"
	"iracema/compile"
	"iracema/interpreter"
	"iracema/parser"
	"os"
)

var (
	disasm = flag.Bool("d", false, "dump instructions")
)

func report(mesg string) {
	fmt.Fprintf(os.Stderr, "\x1b[0;31m%s\x1b[0;0m\n", mesg)
}

func parseFile(file string) *ast.File {
	f, err := os.Open(file)
	if err != nil {
		report(err.Error())
	}

	defer f.Close()

	ast, err := parser.Parse(f)
	if err != nil {
		report(err.Error())
		os.Exit(50)
	}

	return ast
}

func runFile(file string) {
	ast := parseFile(file)

	c := compile.New()
	if *disasm {
		c.Disassemble(ast)
		os.Exit(0)
	}

	meth, err := c.Compile(ast)
	if err != nil {
		report(err.Error())
		os.Exit(60)
	}

	interp := &interpreter.Interpreter{}
	ret, err := interp.Exec(meth)
	if err != nil {
		report(err.Error())
		os.Exit(70)
	}

	fmt.Fprintln(os.Stdout, ret)
	os.Exit(0)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: iracema [flags] [path ...]\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintln(os.Stderr, "error: a .ir file is required")
		os.Exit(2)
		return
	}

	file := flag.Arg(0)
	runFile(file)
}
