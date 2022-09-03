package parser

import (
	"bytes"
	"iracema/ast"
	"testing"
)

func assertError(t *testing.T, code string, expectedErr string) {
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

func setupFunBody(t *testing.T, code string) []ast.Stmt {
	t.Helper()

	buf := new(bytes.Buffer)
	buf.WriteString("fun dummy() {\n")
	buf.WriteString(code)
	buf.WriteString("\n}")
	file, err := Parse(buf)

	if err != nil {
		t.Fatal(err)
	}

	if file == nil {
		t.Fatalf("expected not to be nil")
	}

	if len(file.FunList) != 1 {
		t.Errorf("expected FunList size to be 1, got %d", len(file.FunList))
	}

	return file.FunList[0].Body.Stmts
}

func setupFun(t *testing.T, code string, size int) *ast.FunDecl {
	t.Helper()

	input := bytes.NewBufferString(code)
	file, err := Parse(input)

	if err != nil {
		t.Fatal(err)
	}

	if file == nil {
		t.Fatalf("expected not to be nil")
	}

	if len(file.FunList) != size {
		t.Fatalf("expected FunList size to be %d, got %d", size, len(file.FunList))
	}

	return file.FunList[0]
}

func setupObject(t *testing.T, code string, size int) *ast.ObjectDecl {
	t.Helper()

	input := bytes.NewBufferString(code)
	file, err := Parse(input)

	if err != nil {
		t.Fatal(err)
	}

	if file == nil {
		t.Fatalf("expected not to be nil")
	}

	if len(file.ObjectList) != size {
		t.Fatalf("expected ObjectList to be %d, got %d", size, len(file.ObjectList))
	}

	return file.ObjectList[0]
}

func assertMemberSelector(t *testing.T, expr ast.Expr, expectedBase, expectedMember string) {
	t.Helper()

	mSel, ok := expr.(*ast.MemberSelector)
	if !ok {
		t.Errorf("expected *ast.MemberSelector, got %T", expr)
	}

	assertIdent(t, mSel.Base, expectedBase)
	assertIdent(t, mSel.Member, expectedMember)
}

func assertIdent(t *testing.T, expr ast.Expr, expectedName string) {
	t.Helper()

	ident, ok := expr.(*ast.Ident)
	if !ok {
		t.Errorf("expected *ast.Ident, got %T", expr)
	}

	if ident.Value != expectedName {
		t.Errorf("expected Name to be %q, got %q", expectedName, ident.Value)
	}
}

func assertConst(t *testing.T, expr ast.Expr, expectedName string) {
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

func assertArgumentList(t *testing.T, args []ast.Expr, expectedArgs []string) {
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

func assertLit(t *testing.T, expr ast.Expr, expectedValue string) {
	t.Helper()

	lit, ok := expr.(*ast.BasicLit)
	if !ok {
		t.Errorf("expected to be *ast.BasicLit, got %T", expr)
	}

	if lit.Value != expectedValue {
		t.Errorf("expected *ast.BasicLit.Value to be %q, got %q", expectedValue, lit.Value)
	}
}

type assertParam func(*testing.T, int, *ast.Field)

func assertFunDecl(t *testing.T, stmt ast.Stmt, name string, fn assertParam) *ast.FunDecl {
	t.Helper()

	funDecl, ok := stmt.(*ast.FunDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmt)
	}

	assertIdent(t, funDecl.Name, name)

	for i, field := range funDecl.Parameters {
		fn(t, i, field)
	}

	return funDecl
}
