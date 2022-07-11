package parser

import (
	"iracema/ast"
	"testing"
)

func TestParseForStmt(t *testing.T) {
	code := `for el in elements {}`

	stmts := setupTest(t, code, 1)

	forStmt, ok := stmts[0].(*ast.ForStmt)
	if !ok {
		t.Errorf("expected to be *ast.ForStmt, got %T", stmts[0])
	}

	testIdent(t, forStmt.Element, "el")
	testIdent(t, forStmt.Iterable, "elements")
}
