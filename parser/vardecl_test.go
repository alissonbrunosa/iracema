package parser

import (
	"iracema/ast"
	"testing"
)

func TestVarDecl(t *testing.T) {
	table := []struct {
		scenario       string
		input          string
		wantName       string
		wantType       string
		wantValue      string
		wantParamTypes []string
	}{
		{
			scenario: "without value",
			input:    "var age Int",
			wantName: "age",
			wantType: "Int",
		},
		{
			scenario:  "with value",
			input:     "var age Int = 40",
			wantName:  "age",
			wantType:  "Int",
			wantValue: "40",
		},
		{
			scenario:  "without type with value",
			input:     "var age = 40",
			wantName:  "age",
			wantValue: "40",
		},
		{
			scenario:       "when type has a single argument type",
			input:          "var ages List<Int>",
			wantType:       "List",
			wantName:       "ages",
			wantParamTypes: []string{"Int"},
		},
		{
			scenario:       "when type has two argument types",
			input:          "var cache Map<String, Object>",
			wantType:       "Map",
			wantName:       "cache",
			wantParamTypes: []string{"String", "Object"},
		},
		{
			scenario:       "when type has multi-argument types",
			input:          "var s Something<String, Object, Object>",
			wantType:       "Something",
			wantName:       "s",
			wantParamTypes: []string{"String", "Object", "Object"},
		},
		{
			scenario:       "nested argument type",
			input:          "var list List<List<Int>>",
			wantType:       "List",
			wantName:       "list",
			wantParamTypes: []string{"List"},
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupTest(t, test.input, 1)

			varDecl, ok := stmts[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
			}

			testIdent(t, varDecl.Name, test.wantName)

			if test.wantType != "" {
				testIdent(t, varDecl.Type.Name, test.wantType)

				if len(test.wantParamTypes) != 0 {
					for i, wantParamType := range test.wantParamTypes {
						testConst(t, varDecl.Type.ArgumentTypeList[i], wantParamType)
					}
				}
			}

			if test.wantValue != "" {
				testLit(t, varDecl.Value, test.wantValue)
			}
		})
	}
}
