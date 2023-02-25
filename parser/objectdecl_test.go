package parser

import (
	"iracema/ast"
	"testing"
)

func TestParse_ObjectDecl_(t *testing.T) {
	stmts := setupTest(t, "object Person {}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	if err := assertConstant(objDecl.Name, "Person"); err != nil {
		t.Error(err)
	}
}

func TestParse_ObjectDecl_with_TypeParameters(t *testing.T) {
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

			if err := assertConstant(objDecl.Name, tt.wantName); err != nil {
				t.Error(err)
			}

			for i, paramType := range objDecl.TypeParamList {
				if err := assertConstant(paramType.Name, tt.wantParamTypeNames[i]); err != nil {
					t.Error(err)
				}

				if len(tt.wantParamTypes) != 0 {
					if err := assertConstant(paramType.Type, tt.wantParamTypes[i]); err != nil {
						t.Error(err)
					}
				}
			}
		})
	}
}

func TestParse_ObjectDecl_with_Parent(t *testing.T) {
	stmts := setupTest(t, "object Dog is Animal {}", 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	if err := assertConstant(objDecl.Name, "Dog"); err != nil {
		t.Error(err)
	}
	if err := assertConstant(objDecl.Parent, "Animal"); err != nil {
		t.Error(err)
	}
}

func TestParse_ObjectDecl_withField(t *testing.T) {
	object := `object Person {
  var name String
  const MIN_AGE Int = 18
}`
	stmts := setupTest(t, object, 1)

	objDecl, ok := stmts[0].(*ast.ObjectDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.ObjectDecl, got %T", stmts[0])
	}

	if err := assertConstant(objDecl.Name, "Person"); err != nil {
		t.Error(err)
	}

	for i, field := range objDecl.FieldList {
		if err := assertIdent(field.Name, "name"); err != nil {
			t.Errorf("[FAILED] field name at %d\n\tReason: %s", i, err)
		}

		if err := assertType(field.Type, "String"); err != nil {
			t.Errorf("[FAILED] field type at %d\n\tReason: %s", i, err)
		}
	}

	for i, field := range objDecl.ConstantList {
		if err := assertIdent(field.Name, "MIN_AGE"); err != nil {
			t.Errorf("[FAILED] const name at %d\n\tReason: %s", i, err)
		}

		if err := assertType(field.Type, "Int"); err != nil {
			t.Errorf("[FAILED] const type at %d\n\tReason: %s", i, err)
		}
	}
}
