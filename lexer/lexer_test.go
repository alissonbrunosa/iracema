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
			ExpectedType:    token.EOF,
			ExpectedLiteral: "",
		},
		"ident": {
			Input:           bytes.NewBufferString("name"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "name",
		},
		"ident (ignore whitespace)": {
			Input:           bytes.NewBufferString(" name "),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "name",
		},
		"ident const": {
			Input:           bytes.NewBufferString("Object"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "Object",
		},
		"empty string": {
			Input:           bytes.NewBufferString(`""`),
			ExpectedType:    token.String,
			ExpectedLiteral: "",
		},
		"non empty string": {
			Input:           bytes.NewBufferString(`"test"`),
			ExpectedType:    token.String,
			ExpectedLiteral: "test",
		},
		"block": {
			Input:        bytes.NewBufferString("block"),
			ExpectedType: token.Block,
		},
		"literal true": {
			Input:           bytes.NewBufferString("true"),
			ExpectedType:    token.Bool,
			ExpectedLiteral: "true",
		},
		"literal false": {
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
		"NotEqual": {
			Input:        bytes.NewBufferString("!="),
			ExpectedType: token.NotEqual,
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
			ExpectedType: token.LeftParen,
		},
		"RightParenthesis": {
			Input:        bytes.NewBufferString(")"),
			ExpectedType: token.RightParen,
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
		"keyword object": {
			Input:           bytes.NewBufferString("object"),
			ExpectedType:    token.Object,
			ExpectedLiteral: "object",
		},
		"keyword fun": {
			Input:        bytes.NewBufferString("fun"),
			ExpectedType: token.Fun,
		},
		"keyword catch": {
			Input:        bytes.NewBufferString("catch"),
			ExpectedType: token.Catch,
		},
		"literal int": {
			Input:           bytes.NewBufferString("10"),
			ExpectedType:    token.Int,
			ExpectedLiteral: "10",
		},
		"literal float": {
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
		"keyword if": {
			Input:        bytes.NewBufferString("if"),
			ExpectedType: token.If,
		},
		"keyword else": {
			Input:        bytes.NewBufferString("else"),
			ExpectedType: token.Else,
		},
		"keyword stop": {
			Input:        bytes.NewBufferString("stop"),
			ExpectedType: token.Stop,
		},
		"keyword next": {
			Input:        bytes.NewBufferString("next"),
			ExpectedType: token.Next,
		},
		"keyword for": {
			Input:        bytes.NewBufferString("for"),
			ExpectedType: token.For,
		},
		"keyword in": {
			Input:        bytes.NewBufferString("in"),
			ExpectedType: token.In,
		},
		"keyword while": {
			Input:        bytes.NewBufferString("while"),
			ExpectedType: token.While,
		},
		"Great": {
			Input:        bytes.NewBufferString(">"),
			ExpectedType: token.Great,
		},
		"GreatEqual": {
			Input:        bytes.NewBufferString(">="),
			ExpectedType: token.GreatEqual,
		},
		"Less": {
			Input:        bytes.NewBufferString("<"),
			ExpectedType: token.Less,
		},
		"LessEqual": {
			Input:        bytes.NewBufferString("<="),
			ExpectedType: token.LessEqual,
		},
		"keyword return": {
			Input:        bytes.NewBufferString("return"),
			ExpectedType: token.Return,
		},
		"ident with special char (?)": {
			Input:           bytes.NewBufferString("complete?"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "complete?",
		},
		"ident with special char (!)": {
			Input:           bytes.NewBufferString("boom!"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "boom!",
		},
		"ident camel_case (_)": {
			Input:           bytes.NewBufferString("first_name"),
			ExpectedType:    token.Ident,
			ExpectedLiteral: "first_name",
		},
		"slash": {
			Input:        bytes.NewBufferString("/"),
			ExpectedType: token.Slash,
		},
		"star": {
			Input:        bytes.NewBufferString("*"),
			ExpectedType: token.Star,
		},
		"illegal": {
			Input:        bytes.NewBufferString("%"),
			ExpectedType: token.Illegal,
		},
		"keyword is": {
			Input:        bytes.NewBufferString("is"),
			ExpectedType: token.Is,
		},
		"keyword switch": {
			Input:        bytes.NewBufferString("switch"),
			ExpectedType: token.Switch,
		},
		"keyword case": {
			Input:        bytes.NewBufferString("case"),
			ExpectedType: token.Case,
		},
		"keyword default": {
			Input:        bytes.NewBufferString("default"),
			ExpectedType: token.Default,
		},
		"keyword super": {
			Input:        bytes.NewBufferString("super"),
			ExpectedType: token.Super,
		},
		"Semicolon is new line": {
			Input:        bytes.NewBufferString(";"),
			ExpectedType: token.NewLine,
		},
		"keyword and": {
			Input:        bytes.NewBufferString("and"),
			ExpectedType: token.And,
		},
		"keyword or": {
			Input:        bytes.NewBufferString("or"),
			ExpectedType: token.Or,
		},
		"keyword this": {
			Input:        bytes.NewBufferString("this"),
			ExpectedType: token.This,
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
			ExpectedErr: "string not terminated",
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
			expectedToken: []token.Type{token.Ident, token.Assign, token.Int, token.EOF},
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
				token.EOF,
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
				token.LeftParen,
				token.Int,
				token.Plus,
				token.Int,
				token.RightParen,
				token.NewLine,
				token.EOF,
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
				token.EOF,
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

func TestInvalidEscape(t *testing.T) {
	expectedErr := `unknown escape: \m`
	input := bytes.NewBufferString(`"test\m"`)
	l := New(input, func(_ *token.Position, err string) {
		if err != expectedErr {
			t.Errorf("expected error to be %q, got %q", expectedErr, err)
		}
	})

	if tok := l.NextToken(); tok.Type != token.String {
		t.Errorf("expected token to be %q, got %q", token.String, tok.Type)
	}
}
