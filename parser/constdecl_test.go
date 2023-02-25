package parser

import (
	"iracema/ast"
	"testing"
)

func TestParse_ConstDecl_SimpleType(t *testing.T) {
	stmts := setupTest(t, "const PI Float = 3.1415", 1)

	varDecl, ok := stmts[0].(*ast.ConstDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ConstDecl, got %T", stmts[0])
	}

	if err := assertConstant(varDecl.Name, "PI"); err != nil {
		t.Error(err)
	}

	if err := assertType(varDecl.Type, "Float"); err != nil {
		t.Error(err)
	}

	if err := assertLiteral(varDecl.Value, "3.1415"); err != nil {
		t.Error(err)
	}
}

func TestParse_ConstDecl_ParameterizedType(t *testing.T) {
	stmts := setupTest(t, `const EMPTY Array<T> = "DUMMY"`, 1)

	varDecl, ok := stmts[0].(*ast.ConstDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ConstDecl, got %T", stmts[0])
	}

	if err := assertConstant(varDecl.Name, "EMPTY"); err != nil {
		t.Error(err)
	}

	wantType := &wantType{wantType: "Array", args: []any{"T"}}
	if err := assertType(varDecl.Type, wantType); err != nil {
		t.Error(err)
	}
}
