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

	if err := assertIdent(forStmt.Element, "el"); err != nil {
		t.Error(err)
	}

	if err := assertIdent(forStmt.Iterable, "elements"); err != nil {
		t.Fatal(err)
	}
}
