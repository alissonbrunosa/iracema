package typecheck

import (
	"fmt"
	"iracema/parser"
	"os"
	"testing"
)

func assertErrorInFiles(t *testing.T, filename string, expected []string) {
	t.Helper()

	file, err := os.Open(filename)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	err = Check(fileAST)
	errs := err.(ErrList).Clean()
	if len(errs) != len(expected) {
		for _, e := range errs {
			fmt.Println(e)
		}
		t.Fatalf("expected %d errors, got: %d", len(expected), len(errs))
	}

	for i, e := range errs {
		if e.Error() != expected[i] {
			t.Errorf("at %d pos is expected: %v got: %v", i, expected[i], e.Error())
		}
	}

}

func TestVarDecl(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 4 Col: 17] expected 'Int', found 'String' in declaration",
		"[Lin: 7 Col: 17] expected 'Int', found 'Float' in declaration",
		"[Lin: 10 Col: 19] expected 'Float', found 'Int' in declaration",
		"[Lin: 13 Col: 20] expected 'String', found 'Int' in declaration",
	}

	assertErrorInFiles(t, "testdata/vardecl.ir", expectedErrors)
}

func TestAssignStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 7 Col: 9] expected 'String', found 'Int' in assignment",
		"[Lin: 10 Col: 10] assignment mismatch: 2 variables but 1 values",
		"[Lin: 11 Col: 7] assignment mismatch: 1 variables but 2 values",
		"[Lin: 14 Col: 16] expected 'String', found 'Int' in assignment",
		"[Lin: 17 Col: 5] undefined: z",
	}

	assertErrorInFiles(t, "testdata/assignstmt.ir", expectedErrors)
}

func TestCallExpr(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 5 Col: 19] expected 'Float', found 'Int' in declaration",
		"[Lin: 21 Col: 21] object 'Object' has no method 'do'",
		"[Lin: 31 Col: 20] expected 'Int', found 'String' in argument to do",
		"[Lin: 43 Col: 33] expected 'Int', found 'Float' in argument to init",
		"[Lin: 43 Col: 40] expected 'Float', found 'Int' in argument to init",
		"[Lin: 55 Col: 5] no superclass of 'Animal' has method 'do_something'",
		"[Lin: 67 Col: 11] expected 'String', found 'Int' in argument to eat",
		"[Lin: 72 Col: 19] expected 'Float', found 'Int' in declaration",
		"[Lin: 82 Col: 5] wrong number of arguments (given 1, expected 0)",
		"[Lin: 88 Col: 5] wrong number of arguments (given 1, expected 0)",
	}

	assertErrorInFiles(t, "testdata/callexpr.ir", expectedErrors)
}

func TestFunDecl(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 109 Col: 17] function init can not have return value",
		"[Lin: 39 Col: 12] expected 'Int', found 'String' in return statement",
		"[Lin: 43 Col: 12] expected 'Float', found 'Int' in return statement",
		"[Lin: 48 Col: 3] missing return for function: nine",
		"[Lin: 51 Col: 3] missing return for function: ten",
		"[Lin: 57 Col: 3] missing return for function: eleven",
		"[Lin: 63 Col: 3] missing return for function: twelve",
		"[Lin: 70 Col: 3] missing return for function: thirteen",
		"[Lin: 79 Col: 3] missing return for function: fourteen",
		"[Lin: 92 Col: 12] unexpected return value",
		"[Lin: 96 Col: 22] variable a is already defined in function sixteen",
		"[Lin: 100 Col: 9] variable a is already defined in function seventeen",
		"[Lin: 103 Col: 3] missing return for function: infinity_loop",
	}

	assertErrorInFiles(t, "testdata/fundecl.ir", expectedErrors)
}

func TestBinaryExpr(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 193 Col: 19] expected 'Int', found 'Float' in declaration",
		"[Lin: 199 Col: 10] object 'Bool' do not implement '>' operator",
		"[Lin: 202 Col: 10] object 'Bool' do not implement '>=' operator",
		"[Lin: 205 Col: 10] object 'Bool' do not implement '<' operator",
		"[Lin: 208 Col: 10] object 'Bool' do not implement '<=' operator",
		"[Lin: 217 Col: 20] expected 'Int', found 'String' in declaration",
		"[Lin: 220 Col: 12] object 'String' do not implement '*' operator",
		"[Lin: 223 Col: 12] object 'String' do not implement '/' operator",
		"[Lin: 226 Col: 12] object 'String' do not implement '-' operator",
		"[Lin: 232 Col: 16] object 'Object' do not implement '+' operator",
	}

	assertErrorInFiles(t, "testdata/binaryexpr.ir", expectedErrors)
}

func TestIfStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 3 Col: 11] expected 'Bool', found 'Int' in if statement",
		"[Lin: 47 Col: 8] expected 'Bool', found 'Float' in if statement",
	}

	assertErrorInFiles(t, "testdata/ifstmt.ir", expectedErrors)
}

func TestWhileStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 31 Col: 13] expected 'Bool', found 'Int' in while statement",
		"[Lin: 37 Col: 5] stop statement outside loop",
		"[Lin: 41 Col: 5] next statement outside loop",
		"[Lin: 70 Col: 5] stop statement outside loop",
	}

	assertErrorInFiles(t, "testdata/whilestmt.ir", expectedErrors)
}

func TestSwitchStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 16 Col: 12] duplicate case",
		"[Lin: 4 Col: 12] previous case",
		"[Lin: 37 Col: 10] expected 'String', found 'Int' in argument to ==",
		"[Lin: 45 Col: 19] expected 'Int', found 'Float' in declaration",
		"[Lin: 54 Col: 21] expected 'Float', found 'Int' in declaration",
	}

	assertErrorInFiles(t, "testdata/switchstmt.ir", expectedErrors)
}

func TestMemberSelector(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 18 Col: 10] 'Car' object has no field 'year'",
	}

	assertErrorInFiles(t, "testdata/memberselector.ir", expectedErrors)
}

func TestUnaryExpr(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 17 Col: 12] object 'Object' do not implement '-' unary operator",
		"[Lin: 22 Col: 12] object 'Object' do not implement '+' unary operator",
		"[Lin: 27 Col: 12] object 'Int' do not implement '!' unary operator",
		"[Lin: 32 Col: 12] object 'Float' do not implement '!' unary operator",
		"[Lin: 37 Col: 12] object 'String' do not implement '!' unary operator",
		"[Lin: 42 Col: 12] object 'String' do not implement '+' unary operator",
		"[Lin: 47 Col: 12] object 'String' do not implement '-' unary operator",
		"[Lin: 52 Col: 12] object 'Bool' do not implement '-' unary operator",
		"[Lin: 57 Col: 12] object 'Bool' do not implement '+' unary operator",
	}

	assertErrorInFiles(t, "testdata/unaryexpr.ir", expectedErrors)
}
