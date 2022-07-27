package lexer

import (
	"io"
	"io/ioutil"
	"iracema/token"
)

type ErrorHandler func(*token.Position, string)

type Lexer interface {
	NextToken() *token.Token
}

type lexer struct {
	source       []byte
	char         byte
	offset       int
	readOffset   int
	position     *token.Position
	errorHandler ErrorHandler
	readNewLine  bool
}

func (l *lexer) NextToken() *token.Token {
	l.skipWhitespace()
	position := l.position.Snapshot(l.readOffset)
	l.readNewLine = false

	var literal string
	var kind token.Type

	char := l.char
	if isLetter(char) {
		literal = l.readIdent()
		kind = token.Lookup(literal)

		return token.New(kind, literal, position)
	}

	if isDecimal(char) {
		kind, literal = l.readNumber()

		return token.New(kind, literal, position)
	}

	l.advance()
	switch char {
	case '"':
		kind = token.String
		literal = l.readString()
	case '.':
		kind = token.Dot
	case ':':
		kind = token.Colon
	case '=':
		kind = token.Assign
		if l.char == '=' {
			kind = token.Equal
			l.advance()
		}

	case '>':
		kind = token.Great
		if l.char == '=' {
			kind = token.GreatEqual
			l.advance()
		}

	case '<':
		kind = token.Less
		if l.char == '=' {
			kind = token.LessEqual
			l.advance()
		}

	case '@':
		if isDecimal(l.char) {
			l.errorHandler(position, "instance variable can't start with numbers")
		}

		l.pushBack()
		kind = token.Ident
		literal = l.readIdent()
	case ',':
		kind = token.Comma
	case '(':
		kind = token.LeftParenthesis
	case ')':
		l.readNewLine = true
		kind = token.RightParenthesis
	case '[':
		kind = token.LeftBracket
	case ']':
		l.readNewLine = true
		kind = token.RightBracket
	case '{':
		kind = token.LeftBrace
	case '}':
		l.readNewLine = true
		kind = token.RightBrace
	case '-':
		kind = token.Minus
	case '+':
		kind = token.Plus
	case '/':
		kind = token.Slash
	case '*':
		kind = token.Star
	case '!':
		kind = token.Not
		if l.char == '=' {
			kind = token.NotEqual
			l.advance()
		}
	case '\n', ';':
		kind = token.NewLine
	case 0:
		kind = token.Eof
	default:
		kind = token.Illegal
	}

	return token.New(kind, literal, position)
}

func (l *lexer) advance() {
	if l.char == '\n' {
		l.position.AddLine(l.readOffset)
	}

	if l.readOffset < len(l.source) {
		l.offset = l.readOffset
		l.char = l.source[l.readOffset]
		l.readOffset += 1
	} else {
		l.offset = len(l.source)
		l.char = 0
	}
}

func (l *lexer) pushBack() {
	l.offset -= 1
	l.readOffset -= 1
	l.char = l.source[l.readOffset]
}

func (l *lexer) skipWhitespace() {
	for {
		switch l.char {
		case ' ', '\t', '\r', '\n':
			if l.readNewLine && l.char == '\n' {
				return
			}

			l.advance()

		case '#':
			l.skipComment()
			continue
		default:
			return
		}
	}
}

func (l *lexer) skipComment() {
	l.advance()
	for l.char > 0 && l.char != '\n' {
		l.advance()
	}
}

func isLetter(char byte) bool {
	return 'a' <= char && char <= 'z' || 'A' <= char && char <= 'Z' || char == '_'
}

func isSpecialChar(char byte) bool {
	return char == '?' || char == '!'
}

func isDecimal(char byte) bool {
	return '0' <= char && char <= '9'
}

func (l *lexer) readString() string {
	l.readNewLine = true

	offs := l.offset
	for {
		if l.char == '"' {
			l.advance()
			break
		}

		if l.char == 0 {
			l.errorHandler(l.position.Snapshot(l.offset), "string not terminated")
			break
		}

		if l.char == '\n' {
			l.errorHandler(l.position.Snapshot(l.offset), "new line in string")
			break
		}

		l.advance()
	}

	return string(l.source[offs : l.offset-1])
}

func (l *lexer) readIdent() string {
	l.readNewLine = true

	offs := l.offset
	for isLetter(l.char) || isDecimal(l.char) {
		l.advance()
	}

	if isSpecialChar(l.char) {
		l.advance()
	}

	return string(l.source[offs:l.offset])
}

func (l *lexer) readNumber() (token.Type, string) {
	l.readNewLine = true

	offs := l.offset
	for isDecimal(l.char) {
		l.advance()
	}

	if l.char != '.' {
		return token.Int, string(l.source[offs:l.offset])
	}

	l.advance()
	for isDecimal(l.char) {
		l.advance()
	}

	return token.Float, string(l.source[offs:l.offset])
}

func New(input io.Reader, errHandler ErrorHandler) *lexer {
	bytes, err := ioutil.ReadAll(input)
	if err != nil {
		panic("could not read from input" + err.Error())
	}

	l := &lexer{
		source:       bytes,
		errorHandler: errHandler,
		position:     token.NewPosition(),
	}

	l.advance()
	return l
}
