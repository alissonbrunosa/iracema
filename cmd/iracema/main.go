package main

import (
	"flag"
	"fmt"
	"iracema/compile"
	"iracema/interpreter"
	"iracema/parser"
	"os"
)

func printError(mesg string) {
	fmt.Fprintf(os.Stderr, "\x1b[0;31m%s\x1b[0;0m\n", mesg)
}

func executeFile(file string) {
	f, err := os.Open(file)
	if err != nil {
		printError(err.Error())
	}

	defer f.Close()

	ast, err := parser.Parse(f)
	if err != nil {
		printError(err.Error())
		os.Exit(68)
	}

	c := compile.New()
	meth, err := c.Compile(ast)
	if err != nil {
		printError(err.Error())
		os.Exit(70)
	}

	interp := &interpreter.Interpreter{}
	interp.Init(meth)
	ret, err := interp.Dispatch()
	if err != nil {
		printError(err.Error())
		os.Exit(70)
	}

	fmt.Fprintln(os.Stdout, ret)
	os.Exit(0)
}

func dissamble(file string) {
	f, err := os.Open(file)
	if err != nil {
		printError(err.Error())
	}

	defer f.Close()

	ast, err := parser.Parse(f)
	if err != nil {
		printError(err.Error())
		os.Exit(68)
	}

	var c = compile.New()
	if len(os.Args) >= 2 {
		c.Dissamble(ast)
		os.Exit(0)
	}
}

func main() {
	disasm := flag.Bool("d", false, "disassembled instructions")
	flag.Parse()

	if *disasm {
		if len(os.Args) < 3 {
			printError("a file is expected")
			os.Exit(1)
		}

		file := os.Args[len(os.Args)-1]
		dissamble(file)
	}

	file := os.Args[len(os.Args)-1]
	executeFile(file)

}
