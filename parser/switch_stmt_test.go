package parser

import (
	"iracema/ast"
	"testing"
)

func TestParseSwitchStmt(t *testing.T) {
	code := `switch 10 {
			   case 10: puts(10)
			 }`

	stmts := setupFunBody(t, code)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	assertLit(t, switchStmt.Key, "10")

	for _, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		assetIdent(t, callExpr.Method, "puts")
		assertArgumentList(t, callExpr.Arguments, []string{"10"})
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

	stmts := setupFunBody(t, code)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	assertLit(t, switchStmt.Key, "50")

	for _, cc := range switchStmt.Cases {
		exprStmt, ok := cc.Body.Stmts[0].(*ast.ExprStmt)
		if !ok {
			t.Errorf("expected to be *ast.ExprStmt, got %T", cc.Body.Stmts[0])
		}

		callExpr := exprStmt.Expr.(*ast.CallExpr)
		assetIdent(t, callExpr.Method, "puts")
		assertArgumentList(t, callExpr.Arguments, []string{"10"})
	}

	exprStmt, ok := switchStmt.Default.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", switchStmt.Default.Body.Stmts[0])
	}

	callExpr := exprStmt.Expr.(*ast.CallExpr)
	assetIdent(t, callExpr.Method, "puts")
	assertArgumentList(t, callExpr.Arguments, []string{"Default"})
}

func TestParseSwitchStmt_with_MultipleCases(t *testing.T) {
	code := `switch 50 {
  case 10: puts(10)
  case 20: puts(20)
}`

	stmts := setupFunBody(t, code)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	assertLit(t, switchStmt.Key, "50")

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
		assetIdent(t, callExpr.Method, "puts")
		assertArgumentList(t, callExpr.Arguments, params[i])
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

	stmts := setupFunBody(t, code)

	switchStmt, ok := stmts[0].(*ast.SwitchStmt)
	if !ok {
		t.Errorf("expected to be *ast.SwitchStmt, got %T", stmts[0])
	}

	assertLit(t, switchStmt.Key, "50")

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
		assetIdent(t, callExpr.Method, "puts")
		assertArgumentList(t, callExpr.Arguments, params[i])
	}

	exprStmt, ok := switchStmt.Default.Body.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", switchStmt.Default.Body.Stmts[0])
	}

	callExpr := exprStmt.Expr.(*ast.CallExpr)
	assetIdent(t, callExpr.Method, "puts")
	assertArgumentList(t, callExpr.Arguments, []string{"Default"})

}

func TestSwitch_WithInvalidBlock(t *testing.T) {
	code := `fun dummy {
  switch 10 {
    puts(10)
}`

	assertError(t, code, "[Lin: 3 Col: 5] syntax error: expected case, default or }")
}
