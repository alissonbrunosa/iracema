package compile

import (
	"bytes"
	"iracema/bytecode"
	"iracema/lang"
	"iracema/parser"
	"testing"
)

func compile(code string) *lang.Method {
	input := bytes.NewBufferString(code)
	f, err := parser.Parse(input)
	if err != nil {
		panic(err)
	}

	c := New()
	ins, err := c.Compile(f)
	if err != nil {
		panic(err)
	}

	return ins
}

func TestCompile_BinaryExpr(t *testing.T) {
	matchers := []Match{
		expect(bytecode.Push).withOperand(0).toHaveConstant(1),
		expect(bytecode.Push).withOperand(1).toHaveConstant(2),
		expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("+", 1),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	fun := compile("1 + 2")
	for i, instr := range fun.Instrs() {
		matchers[i].Match(t, instr, fun.Constants())
	}
}

func TestCompile_SimpleExpr(t *testing.T) {
	tests := []struct {
		Scenario string
		Code     string
		Matchs   []Match
	}{
		{
			Scenario: "empty source",
			Code:     "",
			Matchs: []Match{
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile single lit",
			Code:     "true",
			Matchs: []Match{
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile constant access",
			Code:     "Int",
			Matchs: []Match{
				expect(bytecode.GetConstant).withOperand(0).toHaveConstant("Int"),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile grouped expr",
			Code:     "(1 + 2)",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(1),
				expect(bytecode.Push).withOperand(1).toHaveConstant(2),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("+", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "multiple lit",
			Code:     "true\nfalse",
			Matchs: []Match{
				expect(bytecode.PushNone),
				expect(bytecode.Return),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile method call",
			Code:     "method()",
			Matchs: []Match{
				expect(bytecode.PushThis),
				expect(bytecode.CallMethod).withOperand(0).toBeMethodCall("method", 0),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile string assign",
			Code:     `a = "string"`,
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant("string"),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile none assign",
			Code:     "a = none",
			Matchs: []Match{
				expect(bytecode.PushNone),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile bool assign",
			Code:     "a = true",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(true),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile float assign",
			Code:     "a = 3.1415",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(3.1415),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile unary operator not",
			Code:     "a = true\n b = !a",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(true),
				expect(bytecode.SetLocal).withOperand(0),
				expect(bytecode.GetLocal).withOperand(0),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("unot", 0),
				expect(bytecode.SetLocal).withOperand(1),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile index expr assign",
			Code:     "a[1] = 3.1415",
			Matchs: []Match{
				expect(bytecode.PushThis),
				expect(bytecode.CallMethod).withOperand(0).toBeMethodCall("a", 0),
				expect(bytecode.Push).withOperand(1).toHaveConstant(1),
				expect(bytecode.Push).withOperand(2).toHaveConstant(3.1415),
				expect(bytecode.CallMethod).withOperand(3).toBeMethodCall("insert", 2),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile assign after lit",
			Code:     "true\na = 100",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile assign after assign",
			Code:     "a = 100\nb = 200",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(1).toHaveConstant(200),
				expect(bytecode.SetLocal).toHaveOperand(1),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile assign with method call",
			Code:     "val = method()",
			Matchs: []Match{
				expect(bytecode.PushThis),
				expect(bytecode.CallMethod).withOperand(0).toBeMethodCall("method", 0),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile method call with args",
			Code:     "plus(1, 2)",
			Matchs: []Match{
				expect(bytecode.PushThis),
				expect(bytecode.Push).withOperand(0).toHaveConstant(1),
				expect(bytecode.Push).withOperand(1).toHaveConstant(2),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("plus", 2),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile index expr",
			Code:     "list = []\nlist[0]",
			Matchs: []Match{
				expect(bytecode.BuildArray).withOperand(0).toHaveConstant(0),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(0).toHaveConstant(0),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("get", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			fun := compile(test.Code)
			for i, instr := range fun.Instrs() {
				test.Matchs[i].Match(t, instr, fun.Constants())
			}
		})
	}
}

func TestCompileIfStmt(t *testing.T) {
	tests := []struct {
		Scenario string
		Code     string
		Matchs   []Match
	}{
		{
			Scenario: "compile empty if stmt",
			Code:     "if 20 > 10 {}",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(4),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile if with body",
			Code:     "if 20 > 10 { a = 100 }",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile if with return",
			Code:     "if 20 < 10 { return 100 }",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("<", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.Return),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile if followed by another stmt",
			Code:     "if 20 == 10 { a = 100 }\n1 * 1",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(4).toHaveConstant(1),
				expect(bytecode.Push).withOperand(5).toHaveConstant(1),
				expect(bytecode.CallMethod).withOperand(6).toBeMethodCall("*", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile if with else",
			Code:     "if 20 > 10 { x = 100 } else { y = 200 }",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(7),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(9),
				expect(bytecode.Push).withOperand(4).toHaveConstant(200),
				expect(bytecode.SetLocal).toHaveOperand(1),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile if with else followed by another stmt",
			Code:     "if 20 <= 10 { a = 100 } else { b = 200 }\n2 - 1",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("<=", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(7),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(9),
				expect(bytecode.Push).withOperand(4).toHaveConstant(200),
				expect(bytecode.SetLocal).toHaveOperand(1),
				expect(bytecode.Push).withOperand(5).toHaveConstant(2),
				expect(bytecode.Push).withOperand(6).toHaveConstant(1),
				expect(bytecode.CallMethod).withOperand(7).toBeMethodCall("-", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},

		{
			Scenario: "compile if with else with return",
			Code:     "if 20 <= 10 { return 100 } else { return 200 }",
			Matchs: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(20),
				expect(bytecode.Push).withOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("<=", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.Push).withOperand(3).toHaveConstant(100),
				expect(bytecode.Return),
				expect(bytecode.Push).withOperand(4).toHaveConstant(200),
				expect(bytecode.Return),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			fun := compile(test.Code)
			for i, instr := range fun.Instrs() {
				test.Matchs[i].Match(t, instr, fun.Constants())
			}
		})
	}
}

func TestCompileWhileStmt(t *testing.T) {
	tests := []struct {
		Scenario string
		Code     string
		Matches  []Match
	}{
		{
			Scenario: "compile while stmt",
			Code:     "while 100 > 100 {}",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(100),
				expect(bytecode.Push).withOperand(1).toHaveConstant(100),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(5),
				expect(bytecode.Jump).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile while with body",
			Code:     "while 140 >= 100 { a = 200 }",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(140),
				expect(bytecode.Push).withOperand(1).toHaveConstant(100),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">=", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(7),
				expect(bytecode.Push).withOperand(3).toHaveConstant(200),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile while with return",
			Code:     "while 140 >= 100 { return 200 }",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(140),
				expect(bytecode.Push).withOperand(1).toHaveConstant(100),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">=", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.Push).withOperand(3).toHaveConstant(200),
				expect(bytecode.Return),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile while with expr",
			Code:     "while 0 == 0 { 20 / 3 }",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(0),
				expect(bytecode.Push).withOperand(1).toHaveConstant(0),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(9),
				expect(bytecode.Push).withOperand(3).toHaveConstant(20),
				expect(bytecode.Push).withOperand(4).toHaveConstant(3),
				expect(bytecode.CallMethod).withOperand(5).toBeMethodCall("/", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile while with next stmt",
			Code:     "while 0 == 0 { next }",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(0),
				expect(bytecode.Push).withOperand(1).toHaveConstant(0),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(5),
				expect(bytecode.Jump).toHaveOperand(0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile while with stop stmt",
			Code: `a = 10
 				   while a >= 0 {
 				     if a == 5 { stop }
 				     a = a - 1
 				   }`,
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(10),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(1).toHaveConstant(0),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">=", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(16),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(3).toHaveConstant(5),
				expect(bytecode.CallMethod).withOperand(4).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(11),
				expect(bytecode.Jump).toHaveOperand(16),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(5).toHaveConstant(1),
				expect(bytecode.CallMethod).withOperand(6).toBeMethodCall("-", 1),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(2),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			fun := compile(test.Code)
			for i, instr := range fun.Instrs() {
				test.Matches[i].Match(t, instr, fun.Constants())
			}
		})
	}
}

func TestCompileForStmt(t *testing.T) {
	tests := []struct {
		Scenario string
		Code     string
		Matches  []Match
	}{
		{
			Scenario: "compile for stmt",
			Code:     "for el in [] {}",
			Matches: []Match{
				expect(bytecode.BuildArray).toHaveOperand(0),
				expect(bytecode.NewIterator),
				expect(bytecode.Iterate),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(2),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile for with body",
			Code:     "for e in [1] { puts(e) }",
			Matches: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(1),
				expect(bytecode.BuildArray).toHaveOperand(1),
				expect(bytecode.NewIterator),
				expect(bytecode.Iterate),
				expect(bytecode.JumpIfFalse).toHaveOperand(11),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.PushThis),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(3),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile for with return",
			Code:     "for el in [] { return 200 }",
			Matches: []Match{
				expect(bytecode.BuildArray).toHaveOperand(0),
				expect(bytecode.NewIterator),
				expect(bytecode.Iterate),
				expect(bytecode.JumpIfFalse).toHaveOperand(8),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Pop),
				expect(bytecode.Push).withOperand(0).toHaveConstant(200),
				expect(bytecode.Return),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile for with next stmt",
			Code:     "for el in [] { next }",
			Matches: []Match{
				expect(bytecode.BuildArray).toHaveOperand(0),
				expect(bytecode.NewIterator),
				expect(bytecode.Iterate),
				expect(bytecode.JumpIfFalse).toHaveOperand(6),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.Jump).toHaveOperand(2),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile for with stop stmt",
			Code: `for el in [] {
				     if el == 2 { stop }
					 puts(el)
				   }`,
			Matches: []Match{
				expect(bytecode.BuildArray).toHaveOperand(0),
				expect(bytecode.NewIterator),
				expect(bytecode.Iterate),
				expect(bytecode.JumpIfFalse).toHaveOperand(16),
				expect(bytecode.SetLocal).toHaveOperand(0),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.Push).withOperand(0).toHaveConstant(2),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(11),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(16),
				expect(bytecode.PushThis),
				expect(bytecode.GetLocal).toHaveOperand(0),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(2),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			fun := compile(test.Code)
			for i, instr := range fun.Instrs() {
				test.Matches[i].Match(t, instr, fun.Constants())
			}
		})
	}
}

func TestCompileSwitchStmt(t *testing.T) {
	tests := []struct {
		Scenario string
		Code     string
		Matches  []Match
	}{
		{
			Scenario: "compile switch stmt",
			Code:     "switch 10 { case 10: puts(10) }",
			Matches: []Match{
				expect(bytecode.Push).toHaveOperand(0).toHaveConstant(10),
				expect(bytecode.Push).toHaveOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(8),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(3).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(4).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile switch stmt multiple cases",
			Code:     "switch 10 { case 10: puts(10); case 20: puts(20) }",
			Matches: []Match{
				expect(bytecode.Push).toHaveOperand(0).toHaveConstant(10),
				expect(bytecode.Push).toHaveOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(9),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(3).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(4).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(17),
				expect(bytecode.Push).toHaveOperand(5).toHaveConstant(10),
				expect(bytecode.Push).toHaveOperand(6).toHaveConstant(20),
				expect(bytecode.CallMethod).withOperand(7).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(17),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(8).toHaveConstant(20),
				expect(bytecode.CallMethod).withOperand(9).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			Scenario: "compile switch stmt full",
			Code:     `switch 10 { case 10: puts(10); case 20: puts(20); default: puts("default") }`,
			Matches: []Match{
				expect(bytecode.Push).toHaveOperand(0).toHaveConstant(10),
				expect(bytecode.Push).toHaveOperand(1).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(9),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(3).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(4).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(22),
				expect(bytecode.Push).toHaveOperand(5).toHaveConstant(10),
				expect(bytecode.Push).toHaveOperand(6).toHaveConstant(20),
				expect(bytecode.CallMethod).withOperand(7).toBeMethodCall("==", 1),
				expect(bytecode.JumpIfFalse).toHaveOperand(18),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(8).toHaveConstant(20),
				expect(bytecode.CallMethod).withOperand(9).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.Jump).toHaveOperand(22),
				expect(bytecode.PushThis),
				expect(bytecode.Push).toHaveOperand(10).toHaveConstant("default"),
				expect(bytecode.CallMethod).withOperand(11).toBeMethodCall("puts", 1),
				expect(bytecode.Pop),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			fun := compile(test.Code)
			for i, instr := range fun.Instrs() {
				test.Matches[i].Match(t, instr, fun.Constants())
			}
		})
	}
}

func TestCompileObjectDecl_Empty(t *testing.T) {
	objMatches := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	checkers := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.DefineObject).toDefine("Person", objMatches),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	code := "object Person {}"

	fun := compile(code)
	for i, instr := range fun.Instrs() {
		checkers[i].Match(t, instr, fun.Constants())
	}
}

func TestCompileObjectDecl_WithFunction(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	objMatches := []Match{
		expect(bytecode.DefineFunction).toDefine("age", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	checkers := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.DefineObject).toDefine("Person", objMatches),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	code := `object Person {
 			   fun age() {}
 			 }`

	fun := compile(code)
	for i, instr := range fun.Instrs() {
		checkers[i].Match(t, instr, fun.Constants())
	}
}

func TestCompileObjectDecl_Empty_WithParent(t *testing.T) {
	objMatches := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	checkers := []Match{
		expect(bytecode.GetConstant).withOperand(0).toHaveConstant("Animal"),
		expect(bytecode.DefineObject).toDefine("Dog", objMatches),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	code := "object Dog is Animal {}"

	fun := compile(code)
	for i, instr := range fun.Instrs() {
		checkers[i].Match(t, instr, fun.Constants())
	}
}

func TestCompileFunDecl(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.Nop),
		expect(bytecode.GetLocal).toHaveOperand(0),
		expect(bytecode.GetLocal).toHaveOperand(1),
		expect(bytecode.CallMethod).withOperand(0).toBeMethodCall("/", 1),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
		expect(bytecode.MatchType).withOperand(1).toHaveConstant("ZeroDivisionError"),
		expect(bytecode.JumpIfFalse).toHaveOperand(16),
		expect(bytecode.SetLocal).toHaveOperand(2),
		expect(bytecode.PushThis),
		expect(bytecode.GetLocal).toHaveOperand(2),
		expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("puts", 1),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
		expect(bytecode.Throw),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("div", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	code := "fun div(a Int, b Int) { a / b } catch(err: ZeroDivisionError) { puts(err) }"

	meth := compile(code)
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompileFunDecl_withCallSuperImplictArgs(t *testing.T) {
	t.Skip("SKIP FOR NOW, PARAMS ARE REQUIRED FOR FUNCTION CALL")

	methMatches := []Match{
		expect(bytecode.PushThis),
		expect(bytecode.GetLocal).toHaveOperand(0),
		expect(bytecode.GetLocal).toHaveOperand(1),
		expect(bytecode.CallSuper),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("do_stuff", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("fun do_stuff(a Int, b Int) { super }")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompileFunDecl_withCallSuperExplicitArgs(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.PushThis),
		expect(bytecode.GetLocal),
		expect(bytecode.Push).withOperand(0).toHaveConstant(10),
		expect(bytecode.CallSuper),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("do_stuff", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("fun do_stuff(a Int, b Int) { super(a, 10) }")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompileFunDecl_withMultipleCatches(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.Nop),
		expect(bytecode.PushThis),
		expect(bytecode.CallMethod).withOperand(0).toBeMethodCall("explode", 0),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
		expect(bytecode.MatchType).withOperand(1).toHaveConstant("Error"),
		expect(bytecode.JumpIfFalse).toHaveOperand(15),
		expect(bytecode.SetLocal).toHaveOperand(0),
		expect(bytecode.PushThis),
		expect(bytecode.GetLocal).toHaveOperand(0),
		expect(bytecode.CallMethod).withOperand(2).toBeMethodCall("puts", 1),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
		expect(bytecode.MatchType).withOperand(3).toHaveConstant("ExplodeError"),
		expect(bytecode.JumpIfFalse).toHaveOperand(24),
		expect(bytecode.SetLocal).toHaveOperand(0),
		expect(bytecode.PushThis),
		expect(bytecode.Push).withOperand(4).toHaveConstant(1),
		expect(bytecode.CallMethod).withOperand(5).toBeMethodCall("exit", 1),
		expect(bytecode.Pop),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
		expect(bytecode.Throw),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("dangerous", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	code := "fun dangerous() { explode() } catch(err: Error) { puts(err) } catch(err: ExplodeError) { exit(1) }"

	meth := compile(code)
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompileFunDecl_withDefaultParams(t *testing.T) {
	funMatch := []Match{
		expect(bytecode.GetLocal).toHaveOperand(0),
		expect(bytecode.JumpIfTrue).toHaveOperand(4),
		expect(bytecode.Push).withOperand(0).toHaveConstant(10),
		expect(bytecode.SetLocal).toHaveOperand(0),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("sum", funMatch),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	fun := compile("fun sum(a Int = 10, b Int) {}")

	for i, instr := range fun.Instrs() {
		top[i].Match(t, instr, fun.Constants())
	}
}

func TestCompileArrayLit(t *testing.T) {
	top := []Match{
		expect(bytecode.Push).withOperand(0).toHaveConstant(1),
		expect(bytecode.Push).withOperand(1).toHaveConstant(2),
		expect(bytecode.BuildArray).toHaveOperand(2),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	fun := compile("[1, 2]")
	for i, instr := range fun.Instrs() {
		top[i].Match(t, instr, fun.Constants())
	}
}

func TestCompileUnaryOperator(t *testing.T) {
	table := []struct {
		code     string
		matchers []Match
	}{
		{
			code: "+10",
			matchers: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("uadd", 0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			code: "-10",
			matchers: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("usub", 0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
		{
			code: "!10",
			matchers: []Match{
				expect(bytecode.Push).withOperand(0).toHaveConstant(10),
				expect(bytecode.CallMethod).withOperand(1).toBeMethodCall("unot", 0),
				expect(bytecode.PushNone),
				expect(bytecode.Return),
			},
		},
	}

	for _, test := range table {
		fun := compile(test.code)
		for i, instr := range fun.Instrs() {
			test.matchers[i].Match(t, instr, fun.Constants())
		}
	}
}

func TestCompileReturnStmt(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.Push).withOperand(0).toHaveConstant(10),
		expect(bytecode.Return),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("do_stuff", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("fun do_stuff() { return 10 }")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}

}

func TestCompileReturnStmt_withoutValue(t *testing.T) {
	methMatches := []Match{
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	top := []Match{
		expect(bytecode.DefineFunction).toDefine("do_stuff", methMatches),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("fun do_stuff() { return }")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompile_LogicalOperator_Or(t *testing.T) {
	top := []Match{
		expect(bytecode.Push).withOperand(0).toHaveConstant(10),
		expect(bytecode.Push).withOperand(1).toHaveConstant(5),
		expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
		expect(bytecode.JumpIfTrue).withOperand(8),
		expect(bytecode.Push).withOperand(3).toHaveConstant(10),
		expect(bytecode.Push).withOperand(4).toHaveConstant(9),
		expect(bytecode.CallMethod).withOperand(5).toBeMethodCall(">", 1),
		expect(bytecode.JumpIfFalse).withOperand(8),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("if 10 > 5 or 10 > 9 {}")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}

func TestCompile_LogicalOperator_And(t *testing.T) {
	top := []Match{
		expect(bytecode.Push).withOperand(0).toHaveConstant(10),
		expect(bytecode.Push).withOperand(1).toHaveConstant(5),
		expect(bytecode.CallMethod).withOperand(2).toBeMethodCall(">", 1),
		expect(bytecode.JumpIfFalse).withOperand(8),
		expect(bytecode.Push).withOperand(3).toHaveConstant(10),
		expect(bytecode.Push).withOperand(4).toHaveConstant(9),
		expect(bytecode.CallMethod).withOperand(5).toBeMethodCall(">", 1),
		expect(bytecode.JumpIfFalse).withOperand(8),
		expect(bytecode.PushNone),
		expect(bytecode.Return),
	}

	meth := compile("if 10 > 5 and 10 > 9 {}")
	for i, instr := range meth.Instrs() {
		top[i].Match(t, instr, meth.Constants())
	}
}
