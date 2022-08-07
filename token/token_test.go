package token

import (
	"testing"
)

func TestLookup(t *testing.T) {
	tests := []struct {
		Ident    string
		Expected Type
	}{
		{
			Ident:    "object",
			Expected: Object,
		},
		{
			Ident:    "true",
			Expected: Bool,
		},
		{
			Ident:    "false",
			Expected: Bool,
		},
		{
			Ident:    "fun",
			Expected: Fun,
		},
		{
			Ident:    "name",
			Expected: Ident,
		},
		{
			Ident:    "a",
			Expected: Ident,
		},
	}

	for _, test := range tests {
		if keyword := Lookup(test.Ident); keyword != test.Expected {
			t.Errorf("expected %q to be %q", keyword, test.Expected)
		}
	}
}

func TestToken_Precedence(t *testing.T) {
	tests := []struct {
		Tok                *Token
		ExpectedPrecedence int
	}{
		{Tok: &Token{Type: Or}, ExpectedPrecedence: 2},
		{Tok: &Token{Type: And}, ExpectedPrecedence: 2},
		{Tok: &Token{Type: Equal}, ExpectedPrecedence: 3},
		{Tok: &Token{Type: Less}, ExpectedPrecedence: 3},
		{Tok: &Token{Type: LessEqual}, ExpectedPrecedence: 3},
		{Tok: &Token{Type: Great}, ExpectedPrecedence: 3},
		{Tok: &Token{Type: GreatEqual}, ExpectedPrecedence: 3},
		{Tok: &Token{Type: Minus}, ExpectedPrecedence: 4},
		{Tok: &Token{Type: Plus}, ExpectedPrecedence: 4},
		{Tok: &Token{Type: Slash}, ExpectedPrecedence: 5},
		{Tok: &Token{Type: Star}, ExpectedPrecedence: 5},
		{Tok: &Token{Type: Return}, ExpectedPrecedence: 0},
	}

	for _, test := range tests {
		if test.Tok.Precedence() != test.ExpectedPrecedence {
			t.Errorf(
				"expected precedence for %q to be %d, got %d",
				test.Tok,
				test.ExpectedPrecedence,
				test.Tok.Precedence(),
			)
		}
	}
}

func TestTokenType_String(t *testing.T) {
	if Star.String() != "*" {
		t.Errorf("expected to be *, got %s", Star.String())
	}
}
