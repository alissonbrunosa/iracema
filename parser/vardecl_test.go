package parser

import (
	"iracema/ast"
	"testing"
)

func TestParse_VarDecl_SimpleType(t *testing.T) {
	stmts := setupTest(t, "var x Int", 1)

	varDecl, ok := stmts[0].(*ast.VarDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
	}

	if err := assertIdent(varDecl.Name, "x"); err != nil {
		t.Error(err)
	}

	if err := assertType(varDecl.Type, "Int"); err != nil {
		t.Error(err)
	}
}

func TestParse_VarDecl_ParameterizedType(t *testing.T) {
	table := []struct {
		scenario string
		input    string
		wantName string
		wantType *wantType
	}{
		{
			scenario: "when type has a single type argument",
			input:    "var ages List<Int>",
			wantName: "ages",
			wantType: &wantType{
				wantType: "List",
				args:     []any{"Int"},
			},
		},
		{
			scenario: "when type has two argument types",
			input:    "var cache Map<String, Object>",
			wantName: "cache",
			wantType: &wantType{
				wantType: "Map",
				args:     []any{"String", "Object"},
			},
		},
		{
			scenario: "when type has multi-argument types",
			input:    "var s Something<String, Object, Int>",
			wantName: "s",
			wantType: &wantType{
				wantType: "Something",
				args:     []any{"String", "Object", "Int"},
			},
		},
		{
			scenario: "nested argument type",
			input:    "var list List<List<Int>>",
			wantName: "list",
			wantType: &wantType{
				wantType: "List",
				args: []any{
					&wantType{wantType: "List", args: []any{"Int"}},
				},
			},
		},
		{
			scenario: "nested with muilt-argument type",
			input:    "var cache Map<String, List<Int>>",
			wantName: "cache",
			wantType: &wantType{
				wantType: "Map",
				args: []any{
					"String",
					&wantType{wantType: "List", args: []any{"Int"}},
				},
			},
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			stmts := setupTest(t, test.input, 1)

			varDecl, ok := stmts[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
			}

			if err := assertIdent(varDecl.Name, test.wantName); err != nil {
				t.Error(err)
			}

			if err := assertType(varDecl.Type, test.wantType); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParse_VarDecl_WithFunctionSignatureType(t *testing.T) {
	table := []struct {
		scenario string
		input    string
		wantType *wantType
	}{
		{
			scenario: "no parameters, no return",
			input:    "var f fun()",
			wantType: new(wantType),
		},
		{
			scenario: "no parameters, with return",
			input:    "var f fun() -> Int",
			wantType: &wantType{returnType: "Int"},
		},
		{
			scenario: "with parameter, no return",
			input:    "var f fun(Int)",
			wantType: &wantType{
				args: []any{"Int"},
			},
		},
		{
			scenario: "with parameter and return",
			input:    "var f fun(Int) -> String",
			wantType: &wantType{
				returnType: "String",
				args:       []any{"Int"},
			},
		},
		{
			scenario: "with parameterized type",
			input:    "var f fun(List<Int>)",
			wantType: &wantType{
				args: []any{
					&wantType{wantType: "List", args: []any{"Int"}},
				},
			},
		},
		{
			scenario: "with parameterized type as return",
			input:    "var f fun() -> List<Int>",
			wantType: &wantType{
				returnType: &wantType{wantType: "List", args: []any{"Int"}},
			},
		},
		{
			scenario: "with function type",
			input:    "var f fun(fun(Float))",
			wantType: &wantType{
				args: []any{
					&wantType{args: []any{"Float"}},
				},
			},
		},
		{
			scenario: "with function type and return function type",
			input:    "var f fun(fun(Int)) -> fun(Int)",
			wantType: &wantType{
				returnType: &wantType{args: []any{"Int"}},
				args: []any{
					&wantType{args: []any{"Int"}},
				},
			},
		},
	}

	for _, row := range table {
		t.Run(row.scenario, func(t *testing.T) {
			stmts := setupTest(t, row.input, 1)

			varDecl, ok := stmts[0].(*ast.VarDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.VarDecl, got %T", stmts[0])
			}

			if err := assertIdent(varDecl.Name, "f"); err != nil {
				t.Error(err)
			}

			if err := assertType(varDecl.Type, row.wantType); err != nil {
				t.Error(err)
			}
		})
	}
}
