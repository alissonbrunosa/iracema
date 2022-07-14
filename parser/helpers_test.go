package parser

import (
	"bytes"
	"iracema/ast"
	"testing"
)

func testParserError(t *testing.T, code string, expectedErr string) {
	t.Helper()

	input := bytes.NewBufferString(code)
	_, err := Parse(input)

	if err == nil {
		t.Fatal("expected to have an error")
	}

	if err.Error() != expectedErr {
		t.Errorf("expected error to be %q, got %q", expectedErr, err.Error())
	}
}

func setupTest(t *testing.T, code string, expectStmts int) []ast.Stmt {
	t.Helper()

	input := bytes.NewBufferString(code)
	file, err := Parse(input)

	if err != nil {
		t.Fatal(err)
	}

	if file == nil {
		t.Error("expected not to be nil")
	}

	if len(file.Stmts) != expectStmts {
		t.Errorf("expected statements to be %d, got %d", expectStmts, len(file.Stmts))
	}

	return file.Stmts
}

func testIdent(t *testing.T, expr ast.Expr, expectedName string) {
	t.Helper()

	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("expected *ast.Ident, got %T", expr)
	}

	if ident.Value != expectedName {
		t.Errorf("expected Name to be %q, got %q", expectedName, ident.Value)
	}
}

func testConst(t *testing.T, expr ast.Expr, expectedName string) {
	t.Helper()

	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("expected *ast.Ident, got %T", expr)
	}

	if !ident.IsConstant() {
		t.Errorf("expected ident to start with up case letter")
	}

	if ident.Value != expectedName {
		t.Errorf("expected Value to be %q, got %q", expectedName, ident.Value)
	}
}

func testArguments(t *testing.T, args []ast.Expr, expectedArgs []string) {
	t.Helper()

	if len(args) != len(expectedArgs) {
		t.Errorf("expected args length to be %d, got %d", len(expectedArgs), len(args))
		return
	}

	for i, arg := range args {
		switch expr := arg.(type) {
		case *ast.BasicLit:
			if expr.Value != expectedArgs[i] {
				t.Errorf("expected arg at %d index to be %q, got %q", i, expectedArgs[i], expr.Value)
			}
		case *ast.BinaryExpr:
			if expr.String() != expectedArgs[i] {
				t.Errorf("expected arg exprto be %q, got %q", expectedArgs[i], expr.String())
			}
		default:
			t.Error("argument type invalid")
		}
	}
}

func testLit(t *testing.T, expr ast.Expr, expectedValue string) {
	t.Helper()

	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		t.Errorf("expected to be *ast.BasicLit, got %T", expr)
	}

	if lit.Value != expectedValue {
		t.Errorf("expected *ast.BasicLit.Value to be %q, got %q", expectedValue, lit.Value)
	}
}
