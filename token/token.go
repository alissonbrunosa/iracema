package token

var keywords = map[string]Type{
	"if":      If,
	"is":      Is,
	"for":     For,
	"switch":  Switch,
	"case":    Case,
	"default": Default,
	"in":      In,
	"stop":    Stop,
	"next":    Next,
	"while":   While,
	"else":    Else,
	"fun":     Fun,
	"none":    None,
	"true":    Bool,
	"false":   Bool,
	"catch":   Catch,
	"block":   Block,
	"object":  Object,
	"return":  Return,
	"super":   Super,
	"or":      Or,
	"and":     And,
	"this":    This,
	"use":     Use,
	"var":     Var,
	"const":   Const,
}

const LowestPrecedence = 0

type Token struct {
	*Position
	Type    Type
	Literal string
}

func (t *Token) Precedence() int {
	switch t.Type {
	case And, Or:
		return 2
	case Equal, NotEqual, Less, LessEqual, Great, GreatEqual:
		return 3
	case Minus, Plus:
		return 4
	case Slash, Star:
		return 5
	}

	return LowestPrecedence
}

func (t *Token) String() string { return t.Type.String() }

func Lookup(ident string) Type {
	if len(ident) == 1 {
		return Ident
	}

	if keyword, ok := keywords[ident]; ok {
		return keyword
	}

	return Ident
}

func New(tokenType Type, literal string, pos *Position) *Token {
	return &Token{
		Type:     tokenType,
		Literal:  literal,
		Position: pos,
	}
}
