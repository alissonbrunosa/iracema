package parser

import (
	"bytes"
	"fmt"
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
		t.Fatalf("expected not to be nil")
	}

	if len(file.Stmts) != expectStmts {
		t.Errorf("expected statements to be %d, got %d", expectStmts, len(file.Stmts))
	}

	return file.Stmts
}

func assertIdent(expr ast.Expr, expectedName string) error {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return fmt.Errorf("expected *ast.Ident, got %T", expr)
	}

	if ident.Value != expectedName {
		return fmt.Errorf("expected Name to be %q, got %q", expectedName, ident.Value)
	}

	return nil
}

func assertConstant(expr ast.Expr, name string) error {
	ident, ok := expr.(*ast.Ident)
	if !ok {
		return fmt.Errorf("expected *ast.Ident, got %T", expr)
	}

	if !ident.IsConstant() {
		return fmt.Errorf("expected ident to start with up case letter")
	}

	if ident.Value != name {
		return fmt.Errorf("expected constant name to be %q, got %q", name, ident.Value)
	}

	return nil
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

func assertLiteral(expr ast.Expr, expectedValue string) error {
	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		return fmt.Errorf("expected to be *ast.BasicLit, got %T", expr)
	}

	if lit.Value != expectedValue {
		return fmt.Errorf("expected *ast.BasicLit.Value to be %q, got %q", expectedValue, lit.Value)
	}

	return nil
}

type assertParam func(*testing.T, int, *ast.VarDecl)

func assertFunDecl(t *testing.T, stmt ast.Stmt, name string, fn assertParam) *ast.FunDecl {
	t.Helper()

	funDecl, ok := stmt.(*ast.FunDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmt)
	}

	if err := assertIdent(funDecl.Name, name); err != nil {
		t.Error(err)
	}

	for i, field := range funDecl.Parameters {
		fn(t, i, field)
	}

	return funDecl
}

func assertType(t ast.Type, want any) error {
	switch _type := t.(type) {
	case *ast.Ident:
		return assertIdent(_type, want.(string))

	case *ast.ParameterizedType:
		return assertParameterizedType(_type, want.(*wantType))

	default:
		return fmt.Errorf("invalid ast.Type: %T", t)
	}
}

type wantType struct {
	wantType string
	args     []any
}

func assertParameterizedType(t *ast.ParameterizedType, wt *wantType) error {
	if err := assertConstant(t.Name, wt.wantType); err != nil {
		return err
	}

	for i, typeArg := range t.TypeArguments {
		if err := assertType(typeArg, wt.args[i]); err != nil {
			return fmt.Errorf("type argument at %d: %s", i, err)
		}
	}

	return nil
}
