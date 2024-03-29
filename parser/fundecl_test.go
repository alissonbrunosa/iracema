package parser

import (
	"iracema/ast"
	"testing"
)

func TestFunDecl(t *testing.T) {
	stmts := setupTest(t, "fun noop() {}", 1)

	funDecl, ok := stmts[0].(*ast.FunDecl)
	funType := funDecl.Type
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
	}

	if err := assertIdent(funType.Name, "noop"); err != nil {
		t.Error(err)
	}

	if len(funType.ParameterList) != 0 {
		t.Errorf("expected paramerer size to be 0, got %d", len(funType.ParameterList))
	}

	if funType.Return != nil {
		t.Errorf("expected Return to be nil")
	}
}

func TestParse_FunDecl_WithParameter(t *testing.T) {
	type wantParam struct {
		wantName string
		wantType string
	}

	table := []struct {
		scenario   string
		input      string
		wantName   string
		wantParams []wantParam
	}{
		{
			scenario: "single params",
			wantName: "println",
			input:    "fun println(o Object) {}",
			wantParams: []wantParam{
				{wantName: "o", wantType: "Object"},
			},
		},
		{
			scenario: "two params",
			wantName: "copy",
			input:    "fun copy(from Object, to Object) {}",
			wantParams: []wantParam{
				{wantName: "from", wantType: "Object"},
				{wantName: "to", wantType: "Object"},
			},
		},
	}

	for _, row := range table {
		t.Run(row.scenario, func(t *testing.T) {
			stmts := setupTest(t, row.input, 1)

			funDecl, ok := stmts[0].(*ast.FunDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
			}

			funType := funDecl.Type
			if err := assertIdent(funType.Name, row.wantName); err != nil {
				t.Fatal(err)
			}

			for i, parameter := range funType.ParameterList {
				wp := row.wantParams[i]
				if err := assertIdent(parameter.Name, wp.wantName); err != nil {
					t.Errorf("param name failed at %d: %s", i, err)
				}

				if err := assertType(parameter.Type, wp.wantType); err != nil {
					t.Errorf("param type failed at %d: %s", i, err)
				}
			}
		})
	}
}

func TestParse_FunDecl_WithParameterizedType(t *testing.T) {
	table := []struct {
		scenario string
		input    string
		wantName string
		wantType *wantType
	}{
		{
			scenario: "parameter with argument type",
			input:    "fun shuffle(l List<Int>) {}",
			wantName: "l",
			wantType: &wantType{
				wantType: "List",
				args:     []any{"Int"},
			},
		},
		{
			scenario: "parameter with two argument types",
			input:    "fun reset(cache Map<String, Object>) {}",
			wantName: "cache",
			wantType: &wantType{
				wantType: "Map",
				args:     []any{"String", "Object"},
			},
		},
		{
			scenario: "parameter with many argument types",
			input:    "fun do_stuff(s Something<String, Object, Object>) {}",
			wantName: "s",
			wantType: &wantType{
				wantType: "Something",
				args:     []any{"String", "Object", "Object"},
			},
		},
		{
			scenario: "parameter with nested argument type",
			input:    "fun flatten(l List<List<Int>>) {}",
			wantName: "l",
			wantType: &wantType{
				wantType: "List",
				args: []any{
					&wantType{wantType: "List", args: []any{"Int"}},
				},
			},
		},
		{
			scenario: "parameter with a normal argument type and a nested one",
			input:    "fun reset(cache Map<String, List<Int>>) {}",
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

	for _, row := range table {
		t.Run(row.scenario, func(t *testing.T) {
			stmts := setupTest(t, row.input, 1)

			funDecl, ok := stmts[0].(*ast.FunDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
			}

			funType := funDecl.Type
			for i, parameter := range funType.ParameterList {
				if err := assertIdent(parameter.Name, row.wantName); err != nil {
					t.Errorf("param name failed at %d: %s", i, err)
				}

				if err := assertType(parameter.Type, row.wantType); err != nil {
					t.Errorf("param type failed at %d: %s", i, err)
				}
			}
		})
	}
}

func TestParse_FunDecl_Return(t *testing.T) {
	table := []struct {
		scenario string
		input    string
		wantType any
	}{
		{
			scenario: "simple return",
			input:    "fun index_of(o Object) -> Int {}",
			wantType: "Int",
		},
		{
			scenario: "return parametized type",
			input:    "fun map(l List<T>) -> List<NT> {}",
			wantType: &wantType{
				wantType: "List",
				args:     []any{"NT"},
			},
		},
		{
			scenario: "return parametized type with nested parametized type",
			input:    "fun group_by_user_id(articles List<Article>) -> Map<Int, List<Article>> {}",
			wantType: &wantType{
				wantType: "Map",
				args: []any{
					"Int",
					&wantType{wantType: "List", args: []any{"Article"}},
				},
			},
		},
	}

	for _, row := range table {
		t.Run(row.scenario, func(t *testing.T) {
			stmts := setupTest(t, row.input, 1)

			funDecl, ok := stmts[0].(*ast.FunDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
			}

			funType := funDecl.Type
			if err := assertType(funType.Return, row.wantType); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestParse_FunDecl_WithFunctionType(t *testing.T) {
	stmts := setupTest(t, "fun handle(fn fun(Int) -> Float) {}", 1)

	funDecl, ok := stmts[0].(*ast.FunDecl)
	if !ok {
		t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
	}

	funType := funDecl.Type
	if err := assertIdent(funType.Name, "handle"); err != nil {
		t.Error(err)
	}

	if len(funType.ParameterList) != 1 {
		t.Fatalf("expected 1 parameter, got %d", len(funType.ParameterList))
	}

	if funType.Return != nil {
		t.Error("expected return to be nil")
	}

	parameter := funType.ParameterList[0]
	if err := assertIdent(parameter.Name, "fn"); err != nil {
		t.Error(err)
	}

	wantType := &wantType{args: []any{"Int"}, returnType: "Float"}
	if err := assertType(parameter.Type, wantType); err != nil {
		t.Error(err)
	}
}

func TestFunDeclWithCatch(t *testing.T) {
	t.Skip("TODO: will get back to this later")
	// tests := []struct {
	// 	Code          string
	// 	ExpectedRef   string
	// 	ExpectedTypes []string
	// }{
	// 	{
	// 		Code:          "fun walk() {} catch(err: Error) {}",
	// 		ExpectedRef:   "err",
	// 		ExpectedTypes: []string{"Error"},
	// 	},
	// 	{
	// 		Code:          "fun walk() {} catch(err: Error) {} catch(err: AnotherError) {}",
	// 		ExpectedRef:   "err",
	// 		ExpectedTypes: []string{"Error", "AnotherError"},
	// 	},
	// }

	// for _, test := range tests {
	// 	stmts := setupTest(t, test.Code, 1)

	// 	funDecl := assertFunDecl(t, stmts[0], "walk", nil)

	// 	for i, catch := range funDecl.Catches {
	// 	}
	// }
}
