package parser

import (
	"iracema/ast"
	"iracema/token"
	"testing"
)

func TestParseIfStmt(t *testing.T) {
	code := `if value == 10 {
			   "Equal"
			 }`

	stmts := setupFunBody(t, code)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	assertIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "10")

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	assertLit(t, exprStmt.Expr, "Equal")
}

func TestParseIfStmtWithElse(t *testing.T) {
	code := `if value == 10 {
			   "Equal"
			 } else {
			   "Not Equal"
			 }`

	stmts := setupFunBody(t, code)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	assertIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "10")

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	assertLit(t, exprStmt.Expr, "Equal")

	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Else)
	}

	exprStmt, ok = elseBlock.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", elseBlock.Stmts[0])
	}
	assertLit(t, exprStmt.Expr, "Not Equal")
}

func TestParseIfStmt_WithElseIfStmt(t *testing.T) {
	code := `if value == 20 {
			   "path 1"
			 } else if value == 30 {
			   "path 2"
			 }`

	stmts := setupFunBody(t, code)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	assertIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "20")

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	assertLit(t, exprStmt.Expr, "path 1")

	ifStmt, ok = ifStmt.Else.(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Else)
	}

	predicate, ok = ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	assertIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "30")

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok = ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	assertLit(t, exprStmt.Expr, "path 2")
}

func TestIfWithInvalidElseBlock(t *testing.T) {
	code := `fun dummy() {
  if value == 20 {
    puts("path 1")
  } else 100
}`

	assertError(t, code, "[Lin: 4 Col: 10] syntax error: expected left brace or if statement")
}
