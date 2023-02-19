package parser

import (
	"iracema/ast"
	"testing"
)

func TestParseObjectDecl(t *testing.T) {
	stmts := setupTest(t, "object Person {}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	testConst(t, objDecl.Name, "Person")
}

func TestParseObjectDecl_with_TypeParameters(t *testing.T) {
	table := []struct {
		scenario           string
		input              string
		wantName           string
		wantParamTypeNames []string
		wantParamTypes     []string
	}{
		{
			scenario:           "single param type",
			input:              "object Array<E> {}",
			wantName:           "Array",
			wantParamTypeNames: []string{"E"},
		},
		{
			scenario:           "multi-param types",
			input:              "object Map<K, V> {}",
			wantName:           "Map",
			wantParamTypeNames: []string{"K", "V"},
		},
		{
			scenario:           "with bounded type",
			input:              "object List<E is Comparable> {}",
			wantName:           "List",
			wantParamTypeNames: []string{"E"},
			wantParamTypes:     []string{"Comparable"},
		},
		{
			scenario:           "with two bounded types",
			input:              "object List<E is Comparable, S is String> {}",
			wantName:           "List",
			wantParamTypeNames: []string{"E", "S"},
			wantParamTypes:     []string{"Comparable", "String"},
		},
	}

	for _, tt := range table {
		t.Run(tt.scenario, func(t *testing.T) {
			stmts := setupTest(t, tt.input, 1)

			objDecl, ok := stmts[0].(*ast.ObjectDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
			}

			testConst(t, objDecl.Name, tt.wantName)

			for i, paramType := range objDecl.ParamTypeList {
				testConst(t, paramType.Name, tt.wantParamTypeNames[i])

				if len(tt.wantParamTypes) != 0 {
					testConst(t, paramType.Type, tt.wantParamTypes[i])
				}
			}
		})
	}
}

func TestParseObjectDecl_with_Parent(t *testing.T) {
	stmts := setupTest(t, "object Dog is Animal {}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	testConst(t, objDecl.Name, "Dog")
	testConst(t, objDecl.Parent, "Animal")
}

func TestParseObjectDecl_withField(t *testing.T) {
	object := `object Person {
  var name String
}`
	stmts := setupTest(t, object, 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	testConst(t, objDecl.Name, "Person")

	for _, field := range objDecl.FieldList {
		testIdent(t, field.Name, "name")
		testConst(t, field.Type.Name, "String")
	}
}
