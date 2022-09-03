package parser

import (
	"iracema/ast"
	"iracema/token"
	"testing"
)

func TestParseAssignStmt(t *testing.T) {
	tests := []struct {
		scenario          string
		input             string
		expectedIdent     []string
		expectedValue     []string
		expectedStmtCount int
	}{
		{
			scenario:          "assign int",
			input:             "a = 10",
			expectedIdent:     []string{"a"},
			expectedValue:     []string{"10"},
			expectedStmtCount: 1,
		},
		{
			scenario:          "assign float",
			input:             "b = 10.10",
			expectedIdent:     []string{"b"},
			expectedValue:     []string{"10.10"},
			expectedStmtCount: 1,
		},
		{
			scenario:          "assign string",
			input:             "x = \"this is string\"",
			expectedIdent:     []string{"x"},
			expectedValue:     []string{"this is string"},
			expectedStmtCount: 1,
		},
		{
			scenario:          "multiple assign ",
			input:             "a, b = 1, 2",
			expectedIdent:     []string{"a", "b"},
			expectedValue:     []string{"1", "2"},
			expectedStmtCount: 1,
		},
	}

	for _, test := range tests {
		tt := test
		t.Run(tt.scenario, func(t *testing.T) {
			stmts := setupFunBody(t, tt.input)

			assignStmt, ok := stmts[0].(*ast.AssignStmt)
			if !ok {
				t.Errorf("Expected to be a *ast.AssignStmt, got %T", stmts[0])
			}

			for i, leftExpr := range assignStmt.Left {
				leftHand, ok := leftExpr.(*ast.Ident)
				if !ok {
					t.Errorf("Expected leftHand to be a *ast.Ident, got %T", leftHand)
				}

				if leftHand.Value != tt.expectedIdent[i] {
					t.Errorf("Expected name in the leftHand to be %q, got %q", tt.expectedValue[i], leftHand.Value)
				}
			}

			for i, rightExpr := range assignStmt.Right {
				rightHand, ok := rightExpr.(*ast.BasicLit)
				if !ok {
					t.Errorf("Expected rightHand to be a *ast.BasicLit, got %T", rightHand)
				}

				if rightHand.Value != tt.expectedValue[i] {
					t.Errorf("Expected value in the rightHand to be %q, got %q", tt.expectedValue, rightHand.Value)
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
			stmts := setupFunBody(t, tt.Code)

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
		input          string
		expectedOutput string
	}{
		{input: "10 - 2", expectedOutput: "(10-2)"},
		{input: "10 + 2", expectedOutput: "(10+2)"},
		{input: "10 * 2", expectedOutput: "(10*2)"},
		{input: "10 / 2", expectedOutput: "(10/2)"},
		{input: "10 / 2 + 5", expectedOutput: "((10/2)+5)"},
		{input: "10 * 2 + 5", expectedOutput: "((10*2)+5)"},
		{input: "10 / 2 * 5", expectedOutput: "((10/2)*5)"},
		{input: "10 + 2 * 5", expectedOutput: "(10+(2*5))"},
		{input: "(10 + 2) * 5", expectedOutput: "((10+2)*5)"},
		{input: "10 / (2 * 5)", expectedOutput: "(10/(2*5))"},
		{input: "!true", expectedOutput: "(!true)"},
		{input: "!!true", expectedOutput: "(!(!true))"},
		{input: "-10 * 10", expectedOutput: "((-10)*10)"},
		{input: "10 + -10 * 10", expectedOutput: "(10+((-10)*10))"},
	}

	for _, test := range tests {
		stmts := setupFunBody(t, test.input)

		stmt := stmts[0]

		output := stmt.String()

		if output != test.expectedOutput {
			t.Errorf("expected output to be %q, got %q\n", test.expectedOutput, output)
		}
	}
}

func TestParseCallExpr(t *testing.T) {
	tests := []struct {
		Code           string
		expectedBase   string
		expectedMember string
		ExpectedArgs   []string
	}{
		{
			Code:           "author.name()",
			expectedBase:   "author",
			expectedMember: "name",
			ExpectedArgs:   []string{},
		},
		{
			Code:           "one.plus(2)",
			expectedBase:   "one",
			expectedMember: "plus",
			ExpectedArgs:   []string{"2"},
		},
		{
			Code:           "math.pow(2, 3)",
			expectedBase:   "math",
			expectedMember: "pow",
			ExpectedArgs:   []string{"2", "3"},
		},
		{
			Code:           "out.println(1 + 2)",
			expectedBase:   "out",
			expectedMember: "println",
			ExpectedArgs:   []string{"(1+2)"},
		},
	}

	for _, test := range tests {
		stmts := setupFunBody(t, test.Code)

		exprStmt, ok := stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected *ast.ExprStmt, got %T", stmts[0])
		}

		callExpr, ok := exprStmt.Expr.(*ast.MethodCallExpr)
		if !ok {
			t.Errorf("expected *ast.CallExpr, got %T", exprStmt.Expr)
		}

		assertMemberSelector(t, callExpr.Selector, test.expectedBase, test.expectedMember)
		assertArgumentList(t, callExpr.Arguments, test.ExpectedArgs)
	}
}

func TestErrorParse(t *testing.T) {
	tests := []struct {
		Scenario    string
		Code        string
		ExpectedErr string
	}{
		{
			Scenario:    "Missing comma in parameter list",
			Code:        `fun name(arg1 Int arg2 Int) {}`,
			ExpectedErr: "[Lin: 1 Col: 19] syntax error: missing , or )",
		},
		{
			Scenario:    "Missing closing brace",
			Code:        "fun name {\n",
			ExpectedErr: "[Lin: 2 Col: 1] syntax error: expected '}', found 'EOF'",
		},
		{
			Scenario:    "var decl without type and value",
			Code:        "fun dummy { var a }",
			ExpectedErr: "[Lin: 1 Col: 19] syntax error: expected 'Ident', found '}'",
		},
		{
			Scenario:    "object declaration without a constant",
			Code:        "object car {}",
			ExpectedErr: "[Lin: 1 Col: 8] syntax error: expected ident to be a constant",
		},
		{
			Scenario:    "when statement is not VarDecl, FunDecl or ObjectDecl",
			Code:        "10 + 10",
			ExpectedErr: "[Lin: 1 Col: 1] syntax error: unexpected Int, expecting VarDecl, FunDecl or ObjectDecl",
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			assertError(t, test.Code, test.ExpectedErr)
		})
	}
}

func TestFunctionCall(t *testing.T) {
	tests := []struct {
		code         string
		expectedName string
		expectedArgs []string
	}{
		{
			code:         "println()",
			expectedName: "println",
			expectedArgs: []string{},
		},
		{
			code:         `println("Hello")`,
			expectedName: "println",
			expectedArgs: []string{"Hello"},
		},
	}

	for _, test := range tests {
		stmts := setupFunBody(t, test.code)

		exprStmt := stmts[0].(*ast.ExprStmt)

		fCall := exprStmt.Expr.(*ast.FunctionCallExpr)
		assertIdent(t, fCall.Name, test.expectedName)
		assertArgumentList(t, fCall.Arguments, test.expectedArgs)
	}
}

func TestParseObjectDecl(t *testing.T) {
	objDecl := setupObject(t, "object Person{}", 1)
	assertConst(t, objDecl.Name, "Person")
}

func TestParseObjectDecl_with_Parent(t *testing.T) {
	objDecl := setupObject(t, "object Dog is Animal {}", 1)
	assertConst(t, objDecl.Name, "Dog")
	assertConst(t, objDecl.Parent, "Animal")
}

func TestParseObjectDecl_withField(t *testing.T) {
	object := `object Person {
  var name String
}`
	objDecl := setupObject(t, object, 1)

	assertConst(t, objDecl.Name, "Person")

	for _, field := range objDecl.FieldList {
		assertIdent(t, field.Type, "String")
		assertIdent(t, field.Name, "name")
	}
}

func TestFunDecl(t *testing.T) {
	type expectParam = struct {
		name  string
		value string
	}

	tests := []struct {
		scenario       string
		code           string
		expectedParams []expectParam
	}{
		{
			scenario: "no params",
			code:     "fun calc {}",
		},
		{
			scenario: "single params",
			code:     "fun calc(a Int) {}",
			expectedParams: []expectParam{
				{name: "a"},
			},
		},
		{
			scenario: "params has default",
			code:     "fun calc(a Int = 1, b Int = 2) {}",
			expectedParams: []expectParam{
				{name: "a", value: "1"},
				{name: "b", value: "2"},
			},
		},
		{
			scenario: "only one param has default",
			code:     "fun calc(a Int, b Int = 10) {}",
			expectedParams: []expectParam{
				{name: "a"},
				{name: "b", value: "10"},
			},
		},
	}

	for _, test := range tests {
		fun := setupFun(t, test.code, 1)

		assert := func(t *testing.T, pos int, field *ast.Field) {
			assertIdent(t, field.Name, test.expectedParams[pos].name)

			if field.Value != nil {
				assertLit(t, field.Value, test.expectedParams[pos].value)
			}
		}

		assertFunDecl(t, fun, "calc", assert)
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
		fun := setupFun(t, test.Code, 1)

		funDecl := assertFunDecl(t, fun, "walk", nil)

		for i, catch := range funDecl.Catches {
			assertIdent(t, catch.Ref, test.ExpectedRef)
			assertConst(t, catch.Type, test.ExpectedTypes[i])
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
		stmts := setupFunBody(t, test.Code)

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

		assertLit(t, expr.Expr, test.ExpectedValue)
	}
}

func TestParseArrayLiteral(t *testing.T) {
	code := "[1, 2, 3]"
	stmts := setupFunBody(t, code)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	lit, ok := exprStmt.Expr.(*ast.ArrayLit)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ArrayLit, got %T", exprStmt.Expr)
	}

	for i, el := range []string{"1", "2", "3"} {
		assertLit(t, lit.Elements[i], el)
	}
}

func TestParseHashLiteral(t *testing.T) {
	code := "{ 1: 10, 2: 20 }"

	stmts := setupFunBody(t, code)

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

	for i, entry := range lit.Entries {
		assertLit(t, entry.Key, expectedKeyValues[i][0])
		assertLit(t, entry.Value, expectedKeyValues[i][1])
	}
}

func TestParseIndexExpr(t *testing.T) {
	code := "value[10]"

	stmts := setupFunBody(t, code)

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	idxExpr, ok := exprStmt.Expr.(*ast.IndexExpr)
	if !ok {
		t.Errorf("expected first stmt to be *ast.IndexExpr, got %T", exprStmt.Expr)
	}

	assertIdent(t, idxExpr.Expr, "value")
	assertLit(t, idxExpr.Index, "10")
}

func TestParseCodeBlock(t *testing.T) {
	code := "block {}"
	stmts := setupFunBody(t, code)

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
			code:     "super()",
			testFun: func(expr *ast.SuperExpr) {
				if len(expr.Arguments) != 0 {
					t.Errorf("expected .Arguments len to be 0, got %d", len(expr.Arguments))
				}
			},
		},
		{
			scenario: "with arguments",
			code:     "super(value)",
			testFun: func(expr *ast.SuperExpr) {
				if len(expr.Arguments) != 1 {
					t.Errorf("expected .Arguments len to be 1, got %d", len(expr.Arguments))
				}

				for _, arg := range expr.Arguments {
					assertIdent(t, arg, "value")
				}
			},
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupFunBody(t, test.code)

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
	stmts := setupFunBody(t, "a = 10; b = 20")

	if len(stmts) != 2 {
		t.Fatalf("expected to have 2 statements, got %d", len(stmts))
	}

	first, ok := stmts[0].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected ast.AssignStmt, got %T", stmts[0])
	}

	assertIdent(t, first.Left[0], "a")
	assertLit(t, first.Right[0], "10")

	second, ok := stmts[1].(*ast.AssignStmt)
	if !ok {
		t.Fatalf("expected ast.AssignStmt, got %T", stmts[1])
	}

	assertIdent(t, second.Left[0], "b")
	assertLit(t, second.Right[0], "20")
}

func TestParseReturnStmt(t *testing.T) {
	stmts := setupFunBody(t, "return 10")

	returnStmt, ok := stmts[0].(*ast.ReturnStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ReturnStmt, got %T", stmts[0])
	}

	assertLit(t, returnStmt.Value, "10")
}

func TestParseReturnStmt_withoutValue(t *testing.T) {
	stmts := setupFunBody(t, "return")

	returnStmt, ok := stmts[0].(*ast.ReturnStmt)
	if !ok {
		t.Errorf("expected first stmt to be *ast.ReturnStmt, got %T", stmts[0])
	}

	if returnStmt.Value != nil {
		t.Errorf("expected .Value to be nil, got %T", returnStmt.Value)
	}
}

func TestMemberSelector(t *testing.T) {
	stmts := setupFunBody(t, "a.name")

	exprStmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ExprStmt, got %T", stmts[0])
	}

	fs, ok := exprStmt.Expr.(*ast.MemberSelector)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FieldSel, got %T", exprStmt.Expr)
	}

	assertIdent(t, fs.Base, "a")
	assertIdent(t, fs.Member, "name")
}

func TestVarDecl(t *testing.T) {
	table := []struct {
		scenario      string
		input         string
		expectedIdent string
		expectedType  string
		expectedValue string
	}{
		{
			scenario:      "without value",
			input:         "var age Int",
			expectedIdent: "age",
			expectedType:  "Int",
		},
		{
			scenario:      "with value",
			input:         "var age Int = 40",
			expectedIdent: "age",
			expectedType:  "Int",
			expectedValue: "40",
		},
		{
			scenario:      "without type with value",
			input:         "var age = 40",
			expectedIdent: "age",
			expectedValue: "40",
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupFunBody(t, test.input)

			vd, ok := stmts[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
			}

			assertIdent(t, vd.Name, test.expectedIdent)

			if vd.Type != nil {
				assertIdent(t, vd.Type, test.expectedType)
			}

			if vd.Value != nil {
				assertLit(t, vd.Value, test.expectedValue)
			}
		})
	}
}

func TestNewExpr(t *testing.T) {
	stmts := setupFunBody(t, "var o = new Object()")

	vd, ok := stmts[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
	}

	assertIdent(t, vd.Name, "o")

	newExpr, ok := vd.Value.(*ast.NewExpr)
	if !ok {
		t.Fatalf("expected *ast.NewExpr, got %T", vd.Value)
	}

	assertConst(t, newExpr.Type, "Object")
}

func TestNewExpr_WithArguments(t *testing.T) {
	stmts := setupFunBody(t, "var o = new Object(10,10)")

	vd := stmts[0].(*ast.VarDecl)
	newExpr := vd.Value.(*ast.NewExpr)

	assertArgumentList(t, newExpr.Arguments, []string{"10", "10"})
}
