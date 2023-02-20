package parser

import (
	"iracema/ast"
	"testing"
)

func TestVarDecl(t *testing.T) {
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
			scenario: "without value",
			input:    "var age Int",
			wantName: "age",
			wantType: &wantType{
				wantType: "Int",
			},
		},
		{
			scenario:  "with value",
			input:     "var age Int = 40",
			wantValue: "40",
			wantName:  "age",
			wantType: &wantType{
				wantType: "Int",
			},
		},
		{
			scenario:  "without type with value",
			input:     "var age = 40",
			wantValue: "40",
			wantName:  "age",
		},
		{
			scenario: "when type has a single argument type",
			input:    "var ages List<Int>",
			wantName: "ages",
			wantType: &wantType{
				wantType: "List",
				args: []*wantType{
					{wantType: "Int"},
				},
			},
		},
		{
			scenario: "when type has two argument types",
			input:    "var cache Map<String, Object>",
			wantName: "cache",
			wantType: &wantType{
				wantType: "Map",
				args: []*wantType{
					{wantType: "String"},
					{wantType: "Object"},
				},
			},
		},
		{
			scenario: "when type has multi-argument types",
			input:    "var s Something<String, Object, Object>",
			wantName: "s",
			wantType: &wantType{
				wantType: "Something",
				args: []*wantType{
					{wantType: "String"},
					{wantType: "Object"},
					{wantType: "Object"},
				},
			},
		},
		{
			scenario: "nested argument type",
			input:    "var list List<List<Int>>",
			wantName: "list",
			wantType: &wantType{
				wantType: "List",
				args: []*wantType{
					{
						wantType: "List",
						args: []*wantType{
							{wantType: "Int"},
						},
					},
				},
			},
		},
		{
			scenario: "nested with muilt-argument type",
			input:    "var cache Map<String, List<Int>>",
			wantName: "cache",
			wantType: &wantType{
				wantType: "Map",
				args: []*wantType{
					{wantType: "String"},
					{
						wantType: "List",
						args: []*wantType{
							{wantType: "Int"},
						},
					},
				},
			},
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

			varDecl, ok := stmts[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
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
