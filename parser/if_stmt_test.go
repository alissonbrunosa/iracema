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

	stmts := setupTest(t, code, 1)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	if err := assertIdent(predicate.Left, "value"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(predicate.Right, "10"); err != nil {
		t.Error(err)
	}

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	if err := assertLiteral(exprStmt.Expr, "Equal"); err != nil {
		t.Error(err)
	}
}

func TestParseIfStmtWithElse(t *testing.T) {
	code := `if value == 10 {
			   "Equal"
			 } else {
			   "Not Equal"
			 }`

	stmts := setupTest(t, code, 1)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	if err := assertIdent(predicate.Left, "value"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(predicate.Right, "10"); err != nil {
		t.Error(err)
	}

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	if err := assertLiteral(exprStmt.Expr, "Equal"); err != nil {
		t.Error(err)
	}

	elseBlock, ok := ifStmt.Else.(*ast.BlockStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Else)
	}

	exprStmt, ok = elseBlock.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", elseBlock.Stmts[0])
	}
	if err := assertLiteral(exprStmt.Expr, "Not Equal"); err != nil {
		t.Error(err)
	}
}

func TestParseIfStmt_WithElseIfStmt(t *testing.T) {
	code := `if value == 20 {
			   "path 1"
			 } else if value == 30 {
			   "path 2"
			 }`

	stmts := setupTest(t, code, 1)

	ifStmt, ok := stmts[0].(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.IfStmt, got %T", stmts[0])
	}

	predicate, ok := ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	if err := assertIdent(predicate.Left, "value"); err != nil {
		t.Error(err)
	}
	if err := assertLiteral(predicate.Right, "20"); err != nil {
		t.Error(err)
	}

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok := ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	if err := assertLiteral(exprStmt.Expr, "path 1"); err != nil {
		t.Error(err)
	}

	ifStmt, ok = ifStmt.Else.(*ast.IfStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Else)
	}

	predicate, ok = ifStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", ifStmt.Cond)
	}

	if err := assertIdent(predicate.Left, "value"); err != nil {
		t.Error(err)
	}

	if err := assertLiteral(predicate.Right, "30"); err != nil {
		t.Error(err)
	}

	if predicate.Operator.Type != token.Equal {
		t.Errorf("expected operator to be token.Equal, got %q", predicate.Operator)
	}

	exprStmt, ok = ifStmt.Then.Stmts[0].(*ast.ExprStmt)
	if !ok {
		t.Errorf("expected to be *ast.ExprStmt, got %T", ifStmt.Then.Stmts[0])
	}

	if err := assertLiteral(exprStmt.Expr, "path 2"); err != nil {
		t.Error(err)
	}
}

func TestIfWithInvalidElseBlock(t *testing.T) {
	code := `if value == 20 { puts("path 1") } else 100`

	testParserError(t, code, "[Lin: 1 Col: 40] syntax error: expected left brace or if statement")
}
