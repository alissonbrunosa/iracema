package parser

import (
	"bytes"
	"iracema/ast"
	"iracema/token"
	"testing"
)

func TestParseAssignStmt(t *testing.T) {
	tests := []struct {
		Scenario          string
		Code              string
		ExpectedIdent     []string
		ExpectedValue     []string
		ExpectedStmtCount int
	}{
		{
			Scenario:          "assign int",
			Code:              "a = 10",
			ExpectedIdent:     []string{"a"},
			ExpectedValue:     []string{"10"},
			ExpectedStmtCount: 1,
		},
		{
			Scenario:          "assign float",
			Code:              "b = 10.10",
			ExpectedIdent:     []string{"b"},
			ExpectedValue:     []string{"10.10"},
			ExpectedStmtCount: 1,
		},
		{
			Scenario:          "assign string",
			Code:              "x = \"this is string\"",
			ExpectedIdent:     []string{"x"},
			ExpectedValue:     []string{"this is string"},
			ExpectedStmtCount: 1,
		},
		{
			Scenario:          "multiple assign ",
			Code:              "a, b = 1, 2",
			ExpectedIdent:     []string{"a", "b"},
			ExpectedValue:     []string{"1", "2"},
			ExpectedStmtCount: 1,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.Scenario, func(t *testing.T) {
			stmts := setupTest(t, tt.Code, tt.ExpectedStmtCount)

			assignStmt, ok := stmts[0].(*ast.AssignStmt)
			if !ok {
				t.Errorf("Expected to be a *ast.AssignStmt, got %T", stmts[0])
			}

			for i, leftExpr := range assignStmt.Left {
				leftHand, ok := leftExpr.(*ast.Ident)
				if !ok {
					t.Errorf("Expected leftHand to be a *ast.Ident, got %T", leftHand)
				}

				if leftHand.Value != tt.ExpectedIdent[i] {
					t.Errorf("Expected name in the leftHand to be %q, got %q", tt.ExpectedValue[i], leftHand.Value)
				}
			}

			for i, rightExpr := range assignStmt.Right {
				rightHand, ok := rightExpr.(*ast.BasicLit)
				if !ok {
					t.Errorf("Expected rightHand to be a *ast.BasicLit, got %T", rightHand)
				}

				if rightHand.Value != tt.ExpectedValue[i] {
					t.Errorf("Expected value in the rightHand to be %q, got %q", tt.ExpectedValue, rightHand.Value)
				}
			}
		})
	}

}

func TestBinaryExpr(t *testing.T) {
	tests := []struct {
		Scenario           string
		Code               string
		ExpectedLeftValue  string
		ExpectedOperator   token.Type
		ExpectedRightValue string
	}{
		{
			Scenario:           "operation without space between operands",
			Code:               "10+2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Plus,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "add operation",
			Code:               "10 + 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Plus,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "sub operation",
			Code:               "10 - 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Minus,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "div operation",
			Code:               "10 / 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Slash,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "mul operation",
			Code:               "10 * 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Star,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "equal operation",
			Code:               "10 == 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Equal,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "not equal operation",
			Code:               "10 != 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.NotEqual,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "greater than",
			Code:               "10 > 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Great,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "greater or equal than",
			Code:               "10 >= 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.GreatEqual,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "less or equal than",
			Code:               "10 < 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.Less,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "less or equal than",
			Code:               "10 <= 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.LessEqual,
			ExpectedRightValue: "2",
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.Scenario, func(t *testing.T) {
			stmts := setupTest(t, tt.Code, 1)

			exprStmt, ok := stmts[0].(*ast.ExprStmt)
			if !ok {
				t.Errorf("Expected to be a *ast.ExprStmt, got %T", exprStmt)
			}

			expr, ok := exprStmt.Expr.(*ast.BinaryExpr)
			if !ok {
				t.Errorf("Expected to be *ast.BinaryExpr, got %T", expr)
			}

			leftHand, ok := expr.Left.(*ast.BasicLit)
			if !ok {
				t.Errorf("Expected leftHand to be a *ast.BasicLit, got %T", leftHand)
			}

			if leftHand.Value != tt.ExpectedLeftValue {
				t.Errorf("Expected value in the leftHand to be %q, got %q", tt.ExpectedLeftValue, leftHand.Value)
			}

			if expr.Operator.Type != tt.ExpectedOperator {
				t.Errorf("Expected operator to be %q, got %q\n", tt.ExpectedOperator, expr.Operator.Type)
			}

			rightHand, ok := expr.Right.(*ast.BasicLit)
			if !ok {
				t.Errorf("Expected rightHand to be a *ast.BasicLit, got %T", rightHand)
			}

			if rightHand.Value != tt.ExpectedRightValue {
				t.Errorf("Expected value in the rightHand to be %q, got %q", tt.ExpectedRightValue, rightHand.Value)
			}
		})
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		Code           string
		ExpectedOutput string
	}{
		{Code: "10 - 2", ExpectedOutput: "(10-2)"},
		{Code: "10 + 2", ExpectedOutput: "(10+2)"},
		{Code: "10 * 2", ExpectedOutput: "(10*2)"},
		{Code: "10 / 2", ExpectedOutput: "(10/2)"},
		{Code: "10 / 2 + 5", ExpectedOutput: "((10/2)+5)"},
		{Code: "10 * 2 + 5", ExpectedOutput: "((10*2)+5)"},
		{Code: "10 / 2 * 5", ExpectedOutput: "((10/2)*5)"},
		{Code: "10 + 2 * 5", ExpectedOutput: "(10+(2*5))"},
		{Code: "(10 + 2) * 5", ExpectedOutput: "((10+2)*5)"},
		{Code: "10 / (2 * 5)", ExpectedOutput: "(10/(2*5))"},
		{Code: "!true", ExpectedOutput: "(!true)"},
		{Code: "!!true", ExpectedOutput: "(!(!true))"},
		{Code: "-10 * 10", ExpectedOutput: "((-10)*10)"},
		{Code: "10 + -10 * 10", ExpectedOutput: "(10+((-10)*10))"},
	}

	for _, test := range tests {
		input := bytes.NewBufferString(test.Code)
		file, err := Parse(input)

		if err != nil {
			t.Fatal(err)
		}

		output := file.String()

		if output != test.ExpectedOutput {
			t.Errorf("expected output to be %q, got %q\n", test.ExpectedOutput, output)
		}
	}
}

func TestParseCallExpr(t *testing.T) {
	tests := []struct {
		Code             string
		ExpectedReceiver string
		ExpectedMethod   string
		ExpectedArgs     []string
	}{
		{
			Code:             "author.name",
			ExpectedReceiver: "author",
			ExpectedMethod:   "name",
			ExpectedArgs:     []string{},
		},
		{
			Code:             "one.plus(2)",
			ExpectedReceiver: "one",
			ExpectedMethod:   "plus",
			ExpectedArgs:     []string{"2"},
		},
		{
			Code:             "math.pow(2, 3)",
			ExpectedReceiver: "math",
			ExpectedMethod:   "pow",
			ExpectedArgs:     []string{"2", "3"},
		},
		{
			Code:             "out.println(1 + 2)",
			ExpectedReceiver: "out",
			ExpectedMethod:   "println",
			ExpectedArgs:     []string{"(1+2)"},
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		exprStmt, ok := stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected *ast.ExprStmt, got %T", stmts[0])
		}

		callExpr, ok := exprStmt.Expr.(*ast.CallExpr)
		if !ok {
			t.Errorf("expected *ast.CallExpr, got %T", exprStmt.Expr)
		}

		if err := assertIdent(callExpr.Receiver, test.ExpectedReceiver); err != nil {
			t.Error(err)
		}
		if err := assertIdent(callExpr.Method, test.ExpectedMethod); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, test.ExpectedArgs)
	}
}

func TestErrorParse(t *testing.T) {
	tests := []struct {
		scenario  string
		input     string
		wantError string
	}{
		{
			scenario:  "Missing comma in parameter list",
			input:     `fun name(arg1 Int arg2 Int) {}`,
			wantError: "[Lin: 1 Col: 19] syntax error: missing , or )",
		},
		{
			scenario:  "Missing closing brace",
			input:     "fun name() {\n",
			wantError: "[Lin: 2 Col: 1] syntax error: expected '}', found 'EOF'",
		},
		{
			scenario:  "var decl without type and value",
			input:     "var a",
			wantError: "[Lin: 1 Col: 5] syntax error: expected type, found EOF",
		},
		{
			scenario:  "object declaration without a constant",
			input:     "object car {}",
			wantError: "[Lin: 1 Col: 8] syntax error: expected ident to be a constant",
		},
		{
			scenario:  "two statements in the same line",
			input:     "var a = 10 var b = 20",
			wantError: "[Lin: 1 Col: 12] syntax error: unexpected var, expecting EOF or new line",
		},
		{
			scenario:  "two statements in the same line inside function body",
			input:     "fun error() { var a = 10 var b = 20 }",
			wantError: "[Lin: 1 Col: 26] syntax error: unexpected var, expecting } or new line",
		},
		{
			scenario:  "statement other than FunDecl and VarDecl in object's body",
			input:     "object Error { if true {} }",
			wantError: "[Lin: 1 Col: 16] syntax error: unexpected if, expecting FunDecl or VarDecl",
		},
	}

	for _, test := range tests {
		t.Run(test.scenario, func(t *testing.T) {
			testParserError(t, test.input, test.wantError)
		})
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		Code                 string
		ExpectedFunctionName string
		ExpectedArgs         []string
	}{
		{
			Code:                 "println()",
			ExpectedFunctionName: "println",
			ExpectedArgs:         []string{},
		},
		{
			Code:                 `println("Hello")`,
			ExpectedFunctionName: "println",
			ExpectedArgs:         []string{"Hello"},
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		exprStmt := stmts[0].(*ast.ExprStmt)

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		if err := assertIdent(callExpr.Method, test.ExpectedFunctionName); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, test.ExpectedArgs)
	}
}

func TestUnaryExpr(t *testing.T) {
	tests := []struct {
		Code             string
		ExpectedOperator token.Type
		ExpectedValue    string
	}{
		{
			Code:             "!true",
			ExpectedOperator: token.Not,
			ExpectedValue:    "true",
		},

		{
			Code:             "!false",
			ExpectedOperator: token.Not,
			ExpectedValue:    "false",
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		stmtExpr, ok := stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
		}

		expr, ok := stmtExpr.Expr.(*ast.UnaryExpr)
		if !ok {
			t.Errorf("expected first stmt to be *ast.UnaryExpr, got %T", stmtExpr)
		}

		if expr.Operator.Type != test.ExpectedOperator {
			t.Errorf("expected operator to be %q, got %q", test.ExpectedOperator, expr.Operator)
		}

		if err := assertLiteral(expr.Expr, test.ExpectedValue); err != nil {
			t.Error(err)
		}
	}
}

func TestParseArrayLiteral(t *testing.T) {
	code := "[1, 2, 3]"
	stmts := setupTest(t, code, 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	lit, ok := exprStmt.Expr.(*ast.ArrayLit)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ArrayLit, got %T", exprStmt.Expr)
	}

	for i, el := range []string{"1", "2", "3"} {
		if err := assertLiteral(lit.Elements[i], el); err != nil {
			t.Error(err)
		}
	}
}

func TestParseMapLiteral(t *testing.T) {
	code := "{ 1: 10, 2: 20 }"

	stmts := setupTest(t, code, 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	lit, ok := exprStmt.Expr.(*ast.MapLit)
	if !ok {
		t.Errorf("expected first stmt to be *ast.MapLit, got %T", exprStmt.Expr)
	}

	expectedKeyValues := [][]string{
		[]string{"1", "10"},
		[]string{"2", "20"},
	}

	for i, entry := range lit.Entries {
		if err := assertLiteral(entry.Key, expectedKeyValues[i][0]); err != nil {
			t.Error(err)
		}
		if err := assertLiteral(entry.Value, expectedKeyValues[i][1]); err != nil {
			t.Error(err)
		}
	}
}

func TestParseIndexExpr(t *testing.T) {
	code := "value[10]"

	stmts := setupTest(t, code, 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	idxExpr, ok := exprStmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Errorf("expected first stmt to be *ast.IndexExpr, got %T", exprStmt.Expr)
	}

	if err := assertIdent(idxExpr.Expr, "value"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(idxExpr.Index, "10"); err != nil {
		t.Error(err)
	}
}

func TestParseCodeBlock(t *testing.T) {
	code := "block() {}"
	stmts := setupTest(t, code, 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	blockExpr, ok := exprStmt.Expr.(*ast.BlockExpr)
	if !ok {
		t.Errorf("expected first stmt to be *ast.BlockExpr, got %T", exprStmt.Expr)
	}

	if len(blockExpr.Body.Stmts) != 0 {
		t.Errorf("expected block to have 0 stmts")
	}
}

func Test_ParseSuperExpr(t *testing.T) {
	table := []struct {
		scenario string
		code     string
		testFun  func(expr *ast.SuperExpr)
	}{
		{
			scenario: "without args",
			code:     "super",
			testFun: func(expr *ast.SuperExpr) {
				if expr.ExplicitArgs {
					t.Error("expected .ExplictArgs to be false")
				}

				if len(expr.Arguments) != 0 {
					t.Errorf("expected .Arguments len to be 0, got %d", len(expr.Arguments))
				}
			},
		},
		{
			scenario: "explict empty args",
			code:     "super()",
			testFun: func(expr *ast.SuperExpr) {
				if !expr.ExplicitArgs {
					t.Error("expected .ExplicitArgs to be true")
				}

				if len(expr.Arguments) != 0 {
					t.Errorf("expected .Arguments len to be 0, got %d", len(expr.Arguments))
				}
			},
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupTest(t, test.code, 1)

			exprStmt := stmts[0].(*ast.ExprStmt)

			super, ok := exprStmt.Expr.(*ast.SuperExpr)
			if !ok {
				t.Errorf("expected first stmt to be *ast.SuperExpr, got %T", exprStmt.Expr)
			}
			test.testFun(super)
		})
	}
}

func TestParse_withStmtsInTheSameLine(t *testing.T) {
	stmts := setupTest(t, "a = 10; b = 20", 2)

	first, ok := stmts[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected ast.AssignStmt, got %T", stmts[0])
	}

	if err := assertIdent(first.Left[0], "a"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(first.Right[0], "10"); err != nil {
		t.Error(err)
	}

	second, ok := stmts[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected ast.AssignStmt, got %T", stmts[1])
	}

	if err := assertIdent(second.Left[0], "b"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(second.Right[0], "20"); err != nil {
		t.Error(err)
	}
}

func TestParseReturnStmt(t *testing.T) {
	stmts := setupTest(t, "fun do_stuff() { return 10 }", 1)

	funDecl := assertFunDecl(t, stmts[0], "do_stuff", nil)

	returnStmt, ok := funDecl.Body.Stmts[0].(*ast.ReturnStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ReturnStmt, got %T", stmts[0])
	}

	if err := assertLiteral(returnStmt.Value, "10"); err != nil {
		t.Error(err)
	}
}

func TestParseReturnStmt_withoutValue(t *testing.T) {
	stmts := setupTest(t, "fun do_stuff() { return }", 1)

	funDecl := assertFunDecl(t, stmts[0], "do_stuff", nil)

	returnStmt, ok := funDecl.Body.Stmts[0].(*ast.ReturnStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ReturnStmt, got %T", stmts[0])
	}

	if returnStmt.Value != nil {
		t.Errorf("expected .Value to be nil, got %T", returnStmt.Value)
	}
}

func TestFieldSel(t *testing.T) {
	stmts := setupTest(t, "this.name", 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	fs, ok := exprStmt.Expr.(*ast.FieldSel)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FieldSel, got %T", exprStmt.Expr)
	}

	if err := assertIdent(fs.Name, "name"); err != nil {
		t.Error(err)
	}
}
