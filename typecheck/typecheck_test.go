package typecheck

import (
	"bytes"
	"iracema/parser"
	"testing"
)

func assertError(t *testing.T, err error, expected string) {
	t.Helper()

	if err.Error() != expected {
		t.Errorf("\nexpected: %v\n     got: %v", expected, err.Error())
	}

}

func TestVarDecl(t *testing.T) {
	table := []struct {
		scenario      string
		input         string
		expectedError string
	}{
		{
			scenario:      "full variable declaration",
			input:         "var a Int = 10",
			expectedError: "",
		},
		{
			scenario:      "without type",
			input:         "var a = 10",
			expectedError: "",
		},
		{
			scenario:      "without value",
			input:         "var a Int",
			expectedError: "",
		},

		{
			scenario:      "when type and value don't match",
			input:         "var a Int = 3.90",
			expectedError: "cannot use Float as Int value in assignment",
		},
	}

	for _, test := range table {
		in := bytes.NewBufferString(test.input)
		fileAST, err := parser.Parse(in)

		if err != nil {
			t.Errorf("expected no error: %s", err.Error())
		}

		err = Check(fileAST)
		assertError(t, err, test.expectedError)
	}
}

func TestAssignStmt(t *testing.T) {
	table := []struct {
		scenario      string
		input         string
		expectedError string
	}{
		{
			scenario:      "when variable does not exist",
			input:         "x = 10",
			expectedError: "undefined: x",
		},
		{
			scenario:      "assing different type",
			input:         "var x Float; x = 10",
			expectedError: "cannot use Int as Float value in assignment",
		},
	}

	for _, test := range table {
		in := bytes.NewBufferString(test.input)
		fileAST, err := parser.Parse(in)

		if err != nil {
			t.Errorf("expected no error: %s", err.Error())
		}

		err = Check(fileAST)
		assertError(t, err, test.expectedError)
	}
}
