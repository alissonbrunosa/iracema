package parser

import (
	"iracema/ast"
	"iracema/token"
	"testing"
)

func TestParseWhileStmt(t *testing.T) {
	code := `while value > 10 {}`

	stmts := setupFunBody(t, code)

	whileStmt, ok := stmts[0].(*ast.WhileStmt)
	if !ok {
		t.Errorf("expected to be *ast.WhileStmt, got %T", stmts[0])
	}

	predicate, ok := whileStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", whileStmt.Cond)
	}

	assetIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "10")

	if predicate.Operator.Type != token.Great {
		t.Errorf("expected operator to be token.GreaterThan, got %q", predicate.Operator)
	}
}

func TestParseWhileStmtWithStopStmt(t *testing.T) {
	code := "while value > 10 { stop }"

	stmts := setupFunBody(t, code)

	whileStmt, ok := stmts[0].(*ast.WhileStmt)
	if !ok {
		t.Errorf("expected to be *ast.WhileStmt, got %T", stmts[0])
	}

	predicate, ok := whileStmt.Cond.(*ast.BinaryExpr)
	if !ok {
		t.Errorf("expected to be *ast.BinaryExpr, got %T", whileStmt.Cond)
	}

	assetIdent(t, predicate.Left, "value")
	assertLit(t, predicate.Right, "10")

	if predicate.Operator.Type != token.Great {
		t.Errorf("expected operator to be token.GreaterThan, got %q", predicate.Operator)
	}

	_, ok = whileStmt.Body.Stmts[0].(*ast.StopStmt)
	if !ok {
		t.Errorf("expected operator to be *ast.StopStmt., got %T", whileStmt.Body.Stmts[0])
	}
}
