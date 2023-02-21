package parser

import (
	"iracema/ast"
	"testing"
)

func TestConstDecl(t *testing.T) {
	type wantType struct {
		wantType string
		args     []*wantType
	}

	table := []struct {
		scenario  string
		input     string
		wantName  string
		wantType  *wantType
		wantValue string
	}{
		{
			scenario: "without type",
			input:    "const PI = 3.1415",
			wantName: "PI",
		},
		{
			scenario: "with type",
			input:    "const MAX_INT Int = 2147483647",
			wantName: "MAX_INT",
			wantType: &wantType{wantType: "Int"},
		},
	}

	var testParamType func(*testing.T, *ast.Type, *wantType)
	testParamType = func(t *testing.T, tp *ast.Type, wType *wantType) {
		t.Helper()

		testConst(t, tp.Name, wType.wantType)
		for i, arg := range wType.args {
			testParamType(t, tp.ArgumentTypeList[i], arg)
		}
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupTest(t, test.input, 1)

			varDecl, ok := stmts[0].(*ast.ConstDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.ConstDecl, got %T", stmts[0])
			}

			testIdent(t, varDecl.Name, test.wantName)

			if test.wantType != nil {
				testParamType(t, varDecl.Type, test.wantType)
			}

			if test.wantValue != "" {
				testLit(t, varDecl.Value, test.wantValue)
			}
		})
	}
}
