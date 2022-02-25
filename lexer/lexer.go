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
	source      []byte
	currentChar byte
	offset      int
	readOffset  int
	position    *token.Position

	errorHandler ErrorHandler
}

func (l *lexer) NextToken() *token.Token {
	l.skipWhitespace()
	position := l.position.Snapshot(l.readOffset)

	var literal string
	var kind token.Type

	char := l.currentChar
	if isLetter(char) {
		literal = l.readIdent()
		kind = token.Lookup(literal)

		return token.New(kind, literal, position)
	}

	if isDecimal(char) {
		tokenType, literal := l.readNumber()

		return token.New(tokenType, literal, position)
	}

	l.nextChar()
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
		if l.currentChar == '=' {
			kind = token.Equal
			l.nextChar()
		}

	case '>':
		kind = token.GreaterThan
		if l.currentChar == '=' {
			kind = token.GreaterOrEqualThan
			l.nextChar()
		}

	case '<':
		kind = token.LessThan
		if l.currentChar == '=' {
			kind = token.LessOrEqualThan
			l.nextChar()
		}

	case '@':
		if isDecimal(l.currentChar) {
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
		kind = token.RightParenthesis
	case '[':
		kind = token.LeftBracket
	case ']':
		kind = token.RightBracket
	case '{':
		kind = token.LeftBrace
	case '}':
		kind = token.RightBrace
	case '-':
		kind = token.Minus
		if isDecimal(l.currentChar) {
			l.pushBack()
			kind, literal = l.readNumber()
		}
	case '+':
		kind = token.Plus
		if isDecimal(l.currentChar) {
			l.pushBack()
			kind, literal = l.readNumber()
		}
	case '/':
		kind = token.Slash
	case '*':
		kind = token.Star
	case '!':
		kind = token.Not
	case 0:
		kind = token.Eof
	default:
		kind = token.Illegal
	}

	return token.New(kind, literal, position)
}

func (l *lexer) nextChar() {
	if l.readOffset < len(l.source) {
		l.offset = l.readOffset
		l.currentChar = l.source[l.readOffset]
		l.readOffset += 1
	} else {
		l.offset = len(l.source)
		l.currentChar = 0
	}
}

func (l *lexer) pushBack() {
	l.offset -= 1
	l.readOffset -= 1
	l.currentChar = l.source[l.readOffset]
}

func (l *lexer) peek() byte {
	if l.readOffset < len(l.source) {
		return l.source[l.readOffset]
	}
	return 0
}

func (l *lexer) skipWhitespace() {
	for l.currentChar == ' ' || l.currentChar == '\t' || l.currentChar == '\n' || l.currentChar == '\r' {
		if l.currentChar == '\n' {
			l.position.AddLine(l.readOffset)
		}

		l.nextChar()
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
	offs := l.offset

	for {
		if l.currentChar == '"' {
			l.nextChar()
			break
		}

		if l.currentChar == 0 {
			l.errorHandler(l.position.Snapshot(l.offset), "string not terminated")
			break
		}

		if l.currentChar == '\n' {
			l.errorHandler(l.position.Snapshot(l.offset), "new line in string")
			break
		}

		l.nextChar()
	}

	return string(l.source[offs : l.offset-1])
}

func (l *lexer) readIdent() string {
	offs := l.offset
	for isLetter(l.currentChar) || isDecimal(l.currentChar) {
		l.nextChar()
	}

	if isSpecialChar(l.currentChar) {
		l.nextChar()
	}

	return string(l.source[offs:l.offset])
}

func (l *lexer) readNumber() (token.Type, string) {
	offs := l.offset

	for isDecimal(l.currentChar) {
		l.nextChar()
	}

	if l.currentChar != '.' || l.currentChar == '.' && !isDecimal(l.peek()) {
		return token.Int, string(l.source[offs:l.offset])
	}

	l.nextChar()
	for isDecimal(l.currentChar) {
		l.nextChar()
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

	l.nextChar()
	return l
}
