package parser

import (
	"iracema/ast"
	"testing"
)

func TestParseSwitchStmt(t *testing.T) {
	code := `switch 10 {
			   case 10: puts(10)
			 }`

	stmts := setupTest(t, code, 1)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	if err := assertLiteral(switchStmt.Key, "10"); err != nil {
		t.Error(err)
	}

	for _, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		if err := assertIdent(callExpr.Method, "puts"); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, []string{"10"})
	}

	if switchStmt.Default != nil {
		t.Errorf("expected default case to be nil")
	}
}

func TestParseSwitchStmt_with_Default(t *testing.T) {
	code := `switch 50 {
			   case 10: puts(10)
			   default: puts("Default")
			 }`

	stmts := setupTest(t, code, 1)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	if err := assertLiteral(switchStmt.Key, "50"); err != nil {
		t.Error(err)
	}

	for _, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		if err := assertIdent(callExpr.Method, "puts"); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, []string{"10"})
	}

	exprStmt, ok := switchStmt.Default.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", switchStmt.Default.Body.Stmts[0])
	}

	callExpr := exprStmt.Expr.(*ast.CallExpr)
	if err := assertIdent(callExpr.Method, "puts"); err != nil {
		t.Error(err)
	}
	testArguments(t, callExpr.Arguments, []string{"Default"})
}

func TestParseSwitchStmt_with_MultipleCases(t *testing.T) {
	code := `switch 50 {
			   case 10: puts(10)
			   case 20: puts(20)
			 }`

	stmts := setupTest(t, code, 1)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	if err := assertLiteral(switchStmt.Key, "50"); err != nil {
		t.Error(err)
	}

	params := [][]string{
		[]string{"10"},
		[]string{"20"},
	}

	for i, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		if err := assertIdent(callExpr.Method, "puts"); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, params[i])
	}

	if switchStmt.Default != nil {
		t.Errorf("expected default case to be nil")
	}
}

func TestParseSwitchStmt_Full(t *testing.T) {
	code := `switch 50 {
			   case 10: puts(10)
			   case 20: puts(20)
			   default: puts("Default")
			 }`

	stmts := setupTest(t, code, 1)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	if err := assertLiteral(switchStmt.Key, "50"); err != nil {
		t.Error(err)
	}

	params := [][]string{
		[]string{"10"},
		[]string{"20"},
	}

	for i, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		if err := assertIdent(callExpr.Method, "puts"); err != nil {
			t.Error(err)
		}
		testArguments(t, callExpr.Arguments, params[i])
	}

	exprStmt, ok := switchStmt.Default.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", switchStmt.Default.Body.Stmts[0])
	}

	callExpr := exprStmt.Expr.(*ast.CallExpr)
	if err := assertIdent(callExpr.Method, "puts"); err != nil {
		t.Error(err)
	}
	testArguments(t, callExpr.Arguments, []string{"Default"})

}

func TestSwitch_WithInvalidBlock(t *testing.T) {
	code := `switch 10 { puts(10) }`

	testParserError(t, code, "[Lin: 1 Col: 13] syntax error: expected case, default or }")
}
