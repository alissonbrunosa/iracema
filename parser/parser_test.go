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
			Scenario:           "greater than",
			Code:               "10 > 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.GreaterThan,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "greater or equal than",
			Code:               "10 >= 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.GreaterOrEqualThan,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "less or equal than",
			Code:               "10 < 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.LessThan,
			ExpectedRightValue: "2",
		},
		{
			Scenario:           "less or equal than",
			Code:               "10 <= 2",
			ExpectedLeftValue:  "10",
			ExpectedOperator:   token.LessOrEqualThan,
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
		{Code: "-10 * 10", ExpectedOutput: "(-10*10)"},
		{Code: "10 + -10 * 10", ExpectedOutput: "(10+(-10*10))"},
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

		testIdent(t, callExpr.Receiver, test.ExpectedReceiver)
		testIdent(t, callExpr.Method, test.ExpectedMethod)
		testArguments(t, callExpr.Arguments, test.ExpectedArgs)
	}
}

func TestInvalidObjectDecl(t *testing.T) {
	code := `object car {}`

	testParserError(t, code, "[Lin: 1 Col: 8] syntax error: expected ident to be a constant")
}

func TestInvalidFunDecl(t *testing.T) {
	tests := []struct {
		Scenario    string
		Code        string
		ExpectedErr string
	}{

		{
			Scenario:    "Instance Variable as parameter",
			Code:        `fun name(@arg) {}`,
			ExpectedErr: "[Lin: 1 Col: 10] syntax error: argument cannot be an instance variable",
		},
		{
			Scenario:    "Missing comm in parameter list",
			Code:        `fun name(arg1 arg2) {}`,
			ExpectedErr: "[Lin: 1 Col: 15] syntax error: missing ','",
		},
		{
			Scenario:    "Missing closing brace",
			Code:        "fun name {\n",
			ExpectedErr: "[Lin: 2 Col: 1] syntax error: expected '}', found 'Eof'",
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			testParserError(t, test.Code, test.ExpectedErr)
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
		testIdent(t, callExpr.Method, test.ExpectedFunctionName)
		testArguments(t, callExpr.Arguments, test.ExpectedArgs)
	}
}

func TestParseObjectDecl(t *testing.T) {
	stmts := setupTest(t, "object Person{}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	testConst(t, objDecl.Name, "Person")
}

func TestParseObjectDecl_with_Parent(t *testing.T) {
	stmts := setupTest(t, "object Dog is Animal {}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	testConst(t, objDecl.Name, "Dog")
	testConst(t, objDecl.Parent, "Animal")
}

func TestFunDecl(t *testing.T) {
	tests := []struct {
		Code           string
		ExpectedName   string
		ExpectedParams []string
	}{
		{
			Code:         "fun walk {}",
			ExpectedName: "walk",
		},
		{
			Code:           "fun say(msg) {}",
			ExpectedName:   "say",
			ExpectedParams: []string{"msg"},
		},
		{
			Code:           "fun pow(a, b) {}",
			ExpectedName:   "pow",
			ExpectedParams: []string{"a", "b"},
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		funDecl, ok := stmts[0].(*ast.FunDecl)
		if !ok {
			t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
		}

		testIdent(t, funDecl.Name, test.ExpectedName)

		if len(funDecl.Parameters) != len(test.ExpectedParams) {
			t.Fatalf("expected %d params, got %d", len(test.ExpectedParams), len(funDecl.Parameters))
		}

		for i, param := range funDecl.Parameters {
			testIdent(t, param, test.ExpectedParams[i])
		}
	}
}

func TestFunDeclWithCatch(t *testing.T) {
	tests := []struct {
		Code          string
		ExpectedRef   string
		ExpectedTypes []string
	}{
		{
			Code:          "fun walk() {} catch(err: Error) {}",
			ExpectedRef:   "err",
			ExpectedTypes: []string{"Error"},
		},
		{
			Code:          "fun walk() {} catch(err: Error) {} catch(err: AnotherError) {}",
			ExpectedRef:   "err",
			ExpectedTypes: []string{"Error", "AnotherError"},
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		funDecl, ok := stmts[0].(*ast.FunDecl)
		if !ok {
			t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
		}

		for i, catch := range funDecl.Catches {
			testIdent(t, catch.Ref, test.ExpectedRef)
			testConst(t, catch.Type, test.ExpectedTypes[i])
		}
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

		testLit(t, expr.Expr, test.ExpectedValue)
	}
}

func TestParseReturnStmt(t *testing.T) {
	tests := []struct {
		Code         string
		ExpectedExpr string
	}{
		{
			Code:         "return true",
			ExpectedExpr: "true",
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		returnStmt, ok := stmts[0].(*ast.ReturnStmt)
		if !ok {
			t.Errorf("expected first stmt to be *ast.ReturnStmt, got %T", stmts[0])
		}

		testLit(t, returnStmt.Expr, "true")
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
		testLit(t, lit.Elements[i], el)
	}
}

func TestParseHashLiteral(t *testing.T) {
	code := "{ 1: 10, 2: 20 }"

	stmts := setupTest(t, code, 1)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	lit, ok := exprStmt.Expr.(*ast.HashLit)
	if !ok {
		t.Errorf("expected first stmt to be *ast.HashLit, got %T", exprStmt.Expr)
	}

	expectedKeyValues := [][]string{
		[]string{"1", "10"},
		[]string{"2", "20"},
	}

	for i, pair := range lit.Elements {
		testLit(t, pair.Key, expectedKeyValues[i][0])
		testLit(t, pair.Value, expectedKeyValues[i][1])
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

	testIdent(t, idxExpr.Expr, "value")
	testLit(t, idxExpr.Index, "10")
}

func TestParseCodeBlock(t *testing.T) {
	code := "block {}"
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
