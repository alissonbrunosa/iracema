package compile

import (
	"fmt"
	"iracema/ast"
	"iracema/bytecode"
	"iracema/lang"
	"os"
	"strings"
)

func (c *compiler) Disassemble(file *ast.File) {
	_, err := c.Compile(file)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, fragment := range c.fragments {
		header := fragment.name
		rest := (40 - len(header) - 10)
		str := strings.Repeat("=", rest)

		fmt.Println("== disasm:", header, str)
		i := 0
		w := os.Stdout
		for _, ins := range fragment.instrs {
			fmt.Printf("%04d ", i)
			i += 2
			switch ins.opcode {
			case bytecode.Push, bytecode.MatchType, bytecode.GetConstant:
				fmt.Fprintf(w, "%-30s%s\n", ins.opcode, fragment.consts[ins.operand])
			case bytecode.CallMethod, bytecode.CallSuper:
				ci := fragment.consts[ins.operand].(*lang.CallInfo)
				fmt.Fprintf(w, "%-30sname: %s argc:%d\n", ins.opcode, ci.Name(), ci.Argc())
			case bytecode.SetLocal, bytecode.GetLocal:
				fmt.Fprintf(w, "%-30s%s\n", ins.opcode, fragment.locals[ins.operand])
			case bytecode.JumpIfFalse, bytecode.Jump:
				fmt.Fprintf(w, "%-30s%d\n", ins.opcode, ins.operand*2)
			case bytecode.DefineObject:
				m := fragment.consts[ins.operand].(*lang.Method)
				fmt.Fprintf(w, "%-30s%s\n", ins.opcode, m.Name())
			default:
				fmt.Fprintln(w, ins.opcode)
			}
		}
		fmt.Printf("\n")
	}
}
