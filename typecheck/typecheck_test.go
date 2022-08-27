package typecheck

import (
	"iracema/parser"
	"os"
	"testing"
)

func assertError(t *testing.T, err ErrList, expected []string) {
	t.Helper()

	err = err.Clean()

	if len(err) != len(expected) {
		t.Fatalf("expected %d errors, got: %d", len(expected), len(err))
	}

	for i, e := range err {
		if e.Error() != expected[i] {
			t.Errorf("at %d pos is expected: %v got: %v", i, expected[i], e.Error())
		}
	}

}

func TestVarDecl(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 2 Col: 13] cannot use 'String' as 'Int' value in declaration",
		"[Lin: 5 Col: 13] cannot use 'Float' as 'Int' value in declaration",
		"[Lin: 8 Col: 15] cannot use 'Int' as 'Float' value in declaration",
		"[Lin: 11 Col: 16] cannot use 'Int' as 'String' value in declaration",
	}

	file, err := os.Open("testdata/vardecl.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestAssignStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 7 Col: 9] cannot use 'Int' as 'String' value in assignment",
		"[Lin: 10 Col: 10] assignment mismatch: 2 variables but 1 values",
		"[Lin: 11 Col: 7] assignment mismatch: 1 variables but 2 values",
		"[Lin: 14 Col: 16] cannot use 'Int' as 'String' value in assignment",
		"[Lin: 17 Col: 5] undefined: z",
	}

	file, err := os.Open("testdata/assignstmt.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestCallExpr(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 5 Col: 19] cannot use 'Int' as 'Float' value in declaration",
		"[Lin: 21 Col: 21] object 'Object' has no method 'do'",
		"cannot use 'String' as 'Int' in argument to do",
	}

	file, err := os.Open("testdata/callexpr.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestFunDecl(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 39 Col: 12] cannot use 'String' as 'Int' value in return statement",
		"[Lin: 43 Col: 12] cannot use 'Int' as 'Float' value in return statement",
		"[Lin: 48 Col: 3] missing return for function: nine",
		"[Lin: 51 Col: 3] missing return for function: ten",
		"[Lin: 57 Col: 3] missing return for function: eleven",
		"[Lin: 63 Col: 3] missing return for function: twelve",
		"[Lin: 70 Col: 3] missing return for function: thirteen",
		"[Lin: 79 Col: 3] missing return for function: fourteen",
		"[Lin: 92 Col: 12] unexpected return value",
	}

	file, err := os.Open("testdata/fundecl.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestBinaryExpr(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 193 Col: 19] cannot use 'Float' as 'Int' value in declaration",
		"[Lin: 199 Col: 10] object 'Bool' do not implement '>' operator",
		"[Lin: 202 Col: 10] object 'Bool' do not implement '>=' operator",
		"[Lin: 205 Col: 10] object 'Bool' do not implement '<' operator",
		"[Lin: 208 Col: 10] object 'Bool' do not implement '<=' operator",
		"[Lin: 217 Col: 20] cannot use 'String' as 'Int' value in declaration",
		"[Lin: 220 Col: 16] object 'String' do not implement '*' operator",
		"[Lin: 223 Col: 16] object 'String' do not implement '/' operator",
		"[Lin: 226 Col: 16] object 'String' do not implement '-' operator",
		"[Lin: 232 Col: 16] object 'Object' do not implement '+' operator",
	}

	file, err := os.Open("testdata/binaryexpr.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)

}

func TestIfStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 3 Col: 11] expected 'Bool', found 'Int'",
		"[Lin: 47 Col: 8] expected 'Bool', found 'Float'",
	}

	file, err := os.Open("testdata/ifstmt.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestWhileStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 31 Col: 13] expected 'Bool', found 'Int'",
		"[Lin: 37 Col: 5] stop statement outside loop",
		"[Lin: 41 Col: 5] next statement outside loop",
		"[Lin: 70 Col: 5] stop statement outside loop",
	}

	file, err := os.Open("testdata/whilestmt.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}

func TestSwitchStmt(t *testing.T) {
	expectedErrors := []string{
		"[Lin: 16 Col: 12] duplicate case",
		"[Lin: 4 Col: 12] previous case",
		"cannot use 'Int' as 'String' in argument to ==",
		"[Lin: 45 Col: 19] cannot use 'Float' as 'Int' value in declaration",
		"[Lin: 54 Col: 21] cannot use 'Int' as 'Float' value in declaration",
	}

	file, err := os.Open("testdata/switchstmt.ir")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	fileAST, err := parser.Parse(file)
	if err != nil {
		t.Errorf("expected no error: %s", err.Error())
	}

	assertError(t, Check(fileAST), expectedErrors)
}
