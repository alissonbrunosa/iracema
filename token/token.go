package token

type Type int

const (
	Illegal Type = iota
	Eof

	// keywords
	If
	Stop
	Next
	While
	Else
	Fun
	Nil
	Catch
	Block
	Object
	Return

	// literals
	Int
	Float
	String
	Bool

	// arithmetic operations
	Minus
	Plus
	Slash
	Star

	Dot
	Colon
	Not
	Arrow
	Comma
	Assign
	Equal
	LessThan
	LessOrEqualThan
	GreaterThan
	GreaterOrEqualThan
	Ident
	LeftParenthesis
	RightParenthesis
	LeftBracket
	RightBracket
	LeftBrace
	RightBrace
)

var types = [...]string{
	Illegal: "Illegal",
	Eof:     "Eof",

	If:     "if",
	Stop:   "stop",
	Next:   "next",
	While:  "while",
	Else:   "else",
	Fun:    "fun",
	Nil:    "nil",
	Catch:  "catch",
	Block:  "block",
	Object: "object",
	Return: "return",

	Int:    "Int",
	Float:  "Float",
	String: "String",
	Bool:   "Bool",

	Minus: "-",
	Slash: "/",
	Plus:  "+",
	Star:  "*",

	Dot:                ".",
	Colon:              ":",
	Not:                "!",
	Arrow:              "->",
	Comma:              ",",
	Assign:             "=",
	Equal:              "==",
	LessThan:           "<",
	LessOrEqualThan:    "<=",
	GreaterThan:        ">",
	GreaterOrEqualThan: ">=",
	Ident:              "Ident",
	LeftParenthesis:    "(",
	RightParenthesis:   ")",
	LeftBracket:        "[",
	RightBracket:       "]",
	LeftBrace:          "{",
	RightBrace:         "}",
}

var keywords = map[string]Type{
	"if":     If,
	"stop":   Stop,
	"next":   Next,
	"while":  While,
	"else":   Else,
	"fun":    Fun,
	"nil":    Nil,
	"true":   Bool,
	"false":  Bool,
	"catch":  Catch,
	"block":  Block,
	"object": Object,
	"return": Return,
}

const LowestPrecedence = 0

func (t Type) String() string {
	return types[t]
}

type Token struct {
	*Position
	Type    Type
	Literal string
}

func (t *Token) Precedence() int {
	switch t.Type {
	case Equal, LessThan, LessOrEqualThan, GreaterThan, GreaterOrEqualThan:
		return 2
	case Minus, Plus:
		return 3
	case Slash, Star:
		return 4
	}

	return LowestPrecedence
}

func (t *Token) Is(kind Type) bool {
	return t.Type == kind
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
