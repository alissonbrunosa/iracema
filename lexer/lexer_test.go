package lexer

import (
	"bytes"
	"io"
	"iracema/token"
	"testing"
)

func dummyHandler(pos *token.Position, err string) {}

func TestNextToken(t *testing.T) {
	tests := map[string]struct {
		Input           io.Reader
		ExpectedType    token.Type
		ExpectedLiteral string
	}{
		"Eof": {
			Input:           bytes.NewBufferString(""),
			ExpectedType:    token.Eof,
			ExpectedLiteral: "",
		},
		"Ident": {
			Input:           bytes.NewBufferString("name"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "name",
		},
		"Ident Const": {
			Input:           bytes.NewBufferString("Object"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "Object",
		},
		"String": {
			Input:           bytes.NewBufferString(`""`),
			ExpectedType:    token.String,
			ExpectedLiteral: "",
		},
		"block": {
			Input:        bytes.NewBufferString("block"),
			ExpectedType: token.Block,
		},
		"True": {
			Input:           bytes.NewBufferString("true"),
			ExpectedType:    token.Bool,
			ExpectedLiteral: "true",
		},
		"False": {
			Input:           bytes.NewBufferString("false"),
			ExpectedType:    token.Bool,
			ExpectedLiteral: "false",
		},
		"none": {
			Input:        bytes.NewBufferString("none"),
			ExpectedType: token.None,
		},
		"Assign": {
			Input:        bytes.NewBufferString("="),
			ExpectedType: token.Assign,
		},
		"Equal": {
			Input:        bytes.NewBufferString("=="),
			ExpectedType: token.Equal,
		},
		"Not": {
			Input:        bytes.NewBufferString("!"),
			ExpectedType: token.Not,
		},
		"Comma": {
			Input:        bytes.NewBufferString(","),
			ExpectedType: token.Comma,
		},
		"Dot": {
			Input:        bytes.NewBufferString("."),
			ExpectedType: token.Dot,
		},
		"Colon": {
			Input:        bytes.NewBufferString(":"),
			ExpectedType: token.Colon,
		},
		"LeftParenthesis": {
			Input:        bytes.NewBufferString("("),
			ExpectedType: token.LeftParenthesis,
		},
		"RightParenthesis": {
			Input:        bytes.NewBufferString(")"),
			ExpectedType: token.RightParenthesis,
		},
		"LeftBracket": {
			Input:        bytes.NewBufferString("["),
			ExpectedType: token.LeftBracket,
		},
		"RightBracket": {
			Input:        bytes.NewBufferString("]"),
			ExpectedType: token.RightBracket,
		},
		"LeftBrace": {
			Input:        bytes.NewBufferString("{"),
			ExpectedType: token.LeftBrace,
		},
		"RightBrace": {
			Input:        bytes.NewBufferString("}"),
			ExpectedType: token.RightBrace,
		},
		"Object": {
			Input:           bytes.NewBufferString("object"),
			ExpectedType:    token.Object,
			ExpectedLiteral: "object",
		},
		"Fun": {
			Input:        bytes.NewBufferString("fun"),
			ExpectedType: token.Fun,
		},
		"Catch": {
			Input:        bytes.NewBufferString("catch"),
			ExpectedType: token.Catch,
		},
		"Int": {
			Input:           bytes.NewBufferString("10"),
			ExpectedType:    token.Int,
			ExpectedLiteral: "10",
		},
		"Float": {
			Input:           bytes.NewBufferString("10.10"),
			ExpectedType:    token.Float,
			ExpectedLiteral: "10.10",
		},
		"Minus": {
			Input:        bytes.NewBufferString("-"),
			ExpectedType: token.Minus,
		},
		"Plus": {
			Input:        bytes.NewBufferString("+"),
			ExpectedType: token.Plus,
		},
		"If": {
			Input:        bytes.NewBufferString("if"),
			ExpectedType: token.If,
		},
		"Else": {
			Input:        bytes.NewBufferString("else"),
			ExpectedType: token.Else,
		},
		"Stop": {
			Input:        bytes.NewBufferString("stop"),
			ExpectedType: token.Stop,
		},
		"Next": {
			Input:        bytes.NewBufferString("next"),
			ExpectedType: token.Next,
		},
		"for": {
			Input:        bytes.NewBufferString("for"),
			ExpectedType: token.For,
		},
		"in": {
			Input:        bytes.NewBufferString("in"),
			ExpectedType: token.In,
		},
		"While": {
			Input:        bytes.NewBufferString("while"),
			ExpectedType: token.While,
		},
		"GreaterThan": {
			Input:        bytes.NewBufferString(">"),
			ExpectedType: token.Great,
		},
		"GreaterOrEqualThan": {
			Input:        bytes.NewBufferString(">="),
			ExpectedType: token.GreatEqual,
		},
		"LessThan": {
			Input:        bytes.NewBufferString("<"),
			ExpectedType: token.Less,
		},
		"LessOrEqualThan": {
			Input:        bytes.NewBufferString("<="),
			ExpectedType: token.LessEqual,
		},
		"Return": {
			Input:        bytes.NewBufferString("return"),
			ExpectedType: token.Return,
		},
		"Instance variable ident": {
			Input:           bytes.NewBufferString("@name"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "@name",
		},
		"Ident with special char (?)": {
			Input:           bytes.NewBufferString("complete?"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "complete?",
		},
		"Ident with special char (!)": {
			Input:           bytes.NewBufferString("boom!"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "boom!",
		},
		"Ident camel_case (_)": {
			Input:           bytes.NewBufferString("first_name"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "first_name",
		},
		"Slash": {
			Input:        bytes.NewBufferString("/"),
			ExpectedType: token.Slash,
		},
		"Star": {
			Input:        bytes.NewBufferString("*"),
			ExpectedType: token.Star,
		},
		"Illegal": {
			Input:        bytes.NewBufferString("%"),
			ExpectedType: token.Illegal,
		},
		"is": {
			Input:        bytes.NewBufferString("is"),
			ExpectedType: token.Is,
		},
		"switch": {
			Input:        bytes.NewBufferString("switch"),
			ExpectedType: token.Switch,
		},
		"case": {
			Input:        bytes.NewBufferString("case"),
			ExpectedType: token.Case,
		},
		"default": {
			Input:        bytes.NewBufferString("default"),
			ExpectedType: token.Default,
		},
		"super": {
			Input:        bytes.NewBufferString("super"),
			ExpectedType: token.Super,
		},
		"Semicolon is new line": {
			Input:        bytes.NewBufferString(";"),
			ExpectedType: token.NewLine,
		},
	}

	for scenario, test := range tests {
		tt := test
		t.Run(scenario, func(sub *testing.T) {
			l := New(tt.Input, dummyHandler)
			tok := l.NextToken()

			if tok.Type != tt.ExpectedType {
				sub.Fatalf("Expected TokenType to be (%q), got (%q)", tt.ExpectedType, tok.Type)
			}

			if tt.ExpectedLiteral != "" && tok.Literal != tt.ExpectedLiteral {
				sub.Fatalf("Expected Literal to be (%q), got (%q)", tt.ExpectedLiteral, tok.Literal)
			}
		})
	}
}

func TestTokenPositon(t *testing.T) {
	code := `
a = 10
b = "str"
`
	input := bytes.NewBufferString(code)
	l := New(input, dummyHandler)

	expected := []struct {
		Line   int
		Column int
		Kind   token.Type
	}{
		{Line: 2, Column: 1, Kind: token.Ident},
		{Line: 2, Column: 3, Kind: token.Assign},
		{Line: 2, Column: 5, Kind: token.Int},
		{Line: 2, Column: 7, Kind: token.NewLine},
		{Line: 3, Column: 1, Kind: token.Ident},
		{Line: 3, Column: 3, Kind: token.Assign},
		{Line: 3, Column: 5, Kind: token.String},
	}

	for _, expect := range expected {
		tok := l.NextToken()

		if tok.Line() != expect.Line {
			t.Errorf("expected line to be %d, got %d", expect.Line, tok.Line())
		}

		if tok.Column() != expect.Column {
			t.Errorf("expected column to be %d, got %d", expect.Column, tok.Column())
		}

		if tok.Type != expect.Kind {
			t.Errorf("expected kind to be %d, got %d", expect.Kind, tok.Type)
		}
	}
}

func TestInvalidInstanceVariable(t *testing.T) {
	expectedErr := "instance variable can't start with numbers"
	input := bytes.NewBufferString("@1name")
	l := New(input, func(_ *token.Position, err string) {
		if err != expectedErr {
			t.Errorf("expected error to be %q, got %v", expectedErr, err)
		}
	})

	if tok := l.NextToken(); tok.Type != token.Ident {
		t.Errorf("expected token to be %q, got %q", token.Ident, tok.Type)
	}
}

func TestStringError(t *testing.T) {
	tests := []struct {
		Scenario    string
		Input       string
		ExpectedErr string
	}{
		{
			Scenario:    "Non terminated string",
			Input:       `"string`,
			ExpectedErr: "string not terminated",
		},
		{
			Scenario:    "New line in string",
			Input:       "\"string\n",
			ExpectedErr: "new line in string",
		},
	}

	for _, test := range tests {
		t.Run(test.Scenario, func(t *testing.T) {
			input := bytes.NewBufferString(test.Input)
			l := New(input, func(_ *token.Position, err string) {
				if err != test.ExpectedErr {
					t.Errorf("expected error to be %q, got %v", test.ExpectedErr, err)
				}
			})

			if tok := l.NextToken(); tok.Type != token.String {
				t.Errorf("expected token to be %q, got %q", token.String, tok.Type)
			}
		})
	}
}

func TestSkipComment(t *testing.T) {
	table := []struct {
		scenario      string
		source        string
		expectedToken []token.Type
	}{
		{
			scenario:      "same line comment",
			source:        "a = 10 # assing",
			expectedToken: []token.Type{token.Ident, token.Assign, token.Int, token.Eof},
		},
		{
			scenario: "when line above has a comment",
			source:   "# top\na = 10 + 10",
			expectedToken: []token.Type{
				token.Ident,
				token.Assign,
				token.Int,
				token.Plus,
				token.Int,
				token.Eof,
			},
		},
		{
			scenario: "when line below has a comment",
			source:   "a = 2 * (10 + 10)\n# bottom",
			expectedToken: []token.Type{
				token.Ident,
				token.Assign,
				token.Int,
				token.Star,
				token.LeftParenthesis,
				token.Int,
				token.Plus,
				token.Int,
				token.RightParenthesis,
				token.NewLine,
				token.Eof,
			},
		},
		{
			scenario: "surrounded by comments",
			source:   "# top comment\na = 10 # same line\n# bottom comment",
			expectedToken: []token.Type{
				token.Ident,
				token.Assign,
				token.Int,
				token.NewLine,
				token.Eof,
			},
		},
	}

	for _, test := range table {
		t.Run(test.scenario, func(t *testing.T) {
			input := bytes.NewBufferString(test.source)
			l := New(input, nil)

			for i, want := range test.expectedToken {
				got := l.NextToken()

				if got.Type != want {
					t.Errorf("expected  token at %d position to be %s, got %s", i, want, got.Type)
				}
			}
		})
	}
}
