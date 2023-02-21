package parser

import (
	"iracema/ast"
	"testing"
)

func TestFunDecl(t *testing.T) {
	type wantParam = struct {
		name  string
		value string
	}

	tests := []struct {
		scenario   string
		input      string
		wantName   string
		wantParams []wantParam
		wantReturn string
	}{
		{
			scenario: "no params",
			wantName: "calc",
			input:    "fun calc {}",
		},
		{
			scenario: "single params",
			wantName: "calc",
			input:    "fun calc(a Int) {}",
			wantParams: []wantParam{
				{name: "a"},
			},
		},
		{
			scenario: "params has default",
			wantName: "multiply",
			input:    "fun multiply(a Int = 1, b Int = 2) {}",
			wantParams: []wantParam{
				{name: "a", value: "1"},
				{name: "b", value: "2"},
			},
		},
		{
			scenario: "only one param has default",
			wantName: "minus",
			input:    "fun minus(a Int, b Int = 10) {}",
			wantParams: []wantParam{
				{name: "a"},
				{name: "b", value: "10"},
			},
		},
		{
			scenario: "with return type",
			input:    "fun plus(a Int, b Int = 10) -> Int {}",
			wantName: "plus",
			wantParams: []wantParam{
				{name: "a"},
				{name: "b", value: "10"},
			},
			wantReturn: "Int",
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.input, 1)

		funDecl, ok := stmts[0].(*ast.FunDecl)
		if !ok {
			t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
		}

		testIdent(t, funDecl.Name, test.wantName)
		if test.wantReturn != "" {
			testConst(t, funDecl.Return, test.wantReturn)
		}

		for i, parameter := range funDecl.Parameters {
			testIdent(t, parameter.Name, test.wantParams[i].name)

			if parameter.Value != nil {
				testLit(t, parameter.Value, test.wantParams[i].value)
			}
		}

	}
}

func TestFunDeclWithParameterType(t *testing.T) {
	type wantType struct {
		wantType string
		args     []*wantType
	}

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
				args: []*wantType{
					{wantType: "Int"},
				},
			},
		},
		{
			scenario: "parameter with two argument types",
			input:    "fun reset(cache Map<String, Object>) {}",
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
			scenario: "parameter with many argument types",
			input:    "fun do_stuff(s Something<String, Object, Object>) {}",
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
			scenario: "parameter with nested argument type",
			input:    "fun flatten(l List<List<Int>>) {}",
			wantName: "l",
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
			scenario: "parameter with a normal argument type and a nested one",
			input:    "fun reset(cache Map<String, List<Int>>) {}",
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

			funDecl, ok := stmts[0].(*ast.FunDecl)
			if !ok {
				t.Fatalf("expected first stmt to be *ast.FunDecl, got %T", stmts[0])
			}

			for _, parameter := range funDecl.Parameters {
				testIdent(t, parameter.Name, test.wantName)

				testParamType(t, parameter.Type, test.wantType)
			}
		})
	}
}

func TestFunDeclWithCatch(t *testing.T) {
	tests := []struct {
		Code          string
		ExpectedRef   string
		ExpectedTypes []string
	}{
		{
			Code:          "fun walk() {} catch(err: Error) {}",
			ExpectedRef:   "err",
			ExpectedTypes: []string{"Error"},
		},
		{
			Code:          "fun walk() {} catch(err: Error) {} catch(err: AnotherError) {}",
			ExpectedRef:   "err",
			ExpectedTypes: []string{"Error", "AnotherError"},
		},
	}

	for _, test := range tests {
		stmts := setupTest(t, test.Code, 1)

		funDecl := assertFunDecl(t, stmts[0], "walk", nil)

		for i, catch := range funDecl.Catches {
			testIdent(t, catch.Ref, test.ExpectedRef)
			testConst(t, catch.Type, test.ExpectedTypes[i])
		}
	}
}
