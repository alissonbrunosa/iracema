package parser

import (
	"iracema/ast"
	"testing"
)

func TestParseForStmt(t *testing.T) {
	code := `for el in elements {}`
	stmts := setupFunBody(t, code)

	forStmt, ok := stmts[0].(*ast.ForStmt)
	if !ok {
		t.Errorf("expected to be *ast.ForStmt, got %T", stmts[0])
	}

	assertIdent(t, forStmt.Element, "el")
	assertIdent(t, forStmt.Iterable, "elements")
}
