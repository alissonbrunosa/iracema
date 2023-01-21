package lexer

import (
	"io"
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

	if isLetter(l.char) {
		literal := l.readIdent()
		kind := token.Lookup(literal)
		return token.NewToken(kind, literal, position)
	}

	if isDecimal(l.char) {
		kind, literal := l.readNumber()
		return token.NewToken(kind, literal, position)
	}

	switch l.char {
	case '"':
		return token.NewToken(token.String, l.readString(), position)

	case '.':
		l.advance()
		return token.NewToken(token.Dot, "", position)

	case ':':
		l.advance()
		return token.NewToken(token.Colon, "", position)

	case '=':
		l.advance()
		kind := token.Assign

		if l.char == '=' {
			l.advance()
			kind = token.Equal
		}
		return token.NewToken(kind, "", position)

	case '>':
		l.advance()
		l.readNewLine = true
		kind := token.Great

		if l.char == '=' {
			l.advance()
			l.readNewLine = false
			kind = token.GreatEqual
		}

		return token.NewToken(kind, "", position)

	case '<':
		l.advance()
		kind := token.Less

		if l.char == '=' {
			l.advance()
			kind = token.LessEqual
		}
		return token.NewToken(kind, "", position)

	case ',':
		l.advance()
		return token.NewToken(token.Comma, "", position)

	case '(':
		l.advance()
		return token.NewToken(token.LeftParen, "", position)

	case ')':
		l.advance()
		l.readNewLine = true
		return token.NewToken(token.RightParen, "", position)

	case '[':
		l.advance()
		return token.NewToken(token.LeftBracket, "", position)

	case ']':
		l.advance()
		l.readNewLine = true
		return token.NewToken(token.RightBracket, "", position)

	case '{':
		l.advance()
		return token.NewToken(token.LeftBrace, "", position)

	case '}':
		l.advance()
		l.readNewLine = true
		return token.NewToken(token.RightBrace, "", position)

	case '-':
		l.advance()
		kind := token.Minus

		if l.char == '>' {
			l.advance()
			kind = token.Arrow
		}

		return token.NewToken(kind, "", position)

	case '+':
		l.advance()
		return token.NewToken(token.Plus, "", position)

	case '/':
		l.advance()
		return token.NewToken(token.Slash, "", position)

	case '*':
		l.advance()
		return token.NewToken(token.Star, "", position)

	case '!':
		l.advance()
		kind := token.Not

		if l.char == '=' {
			l.advance()
			kind = token.NotEqual
		}
		return token.NewToken(kind, "", position)

	case '\n', ';':
		l.advance()
		return token.NewToken(token.NewLine, "", position)

	case 0:
		return token.NewToken(token.EOF, "", position)

	default:
		l.advance()
		return token.NewToken(token.Illegal, "", position)
	}
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
	l.advance()
	l.readNewLine = true

	start := l.offset
	for {
		if l.char == '"' {
			l.advance()
			break
		}

		if l.char == '\\' {
			l.advance()
			if !l.escape() {
				break
			}
		}

		if l.char <= 0 || l.char == '\n' {
			l.errorHandler(l.position.Snapshot(l.offset), "string not terminated")
			break
		}

		l.advance()
	}

	return string(l.source[start : l.offset-1])
}

func (l *lexer) escape() bool {
	switch l.char {
	case 'a', 'b', 'f', 'n', 'r', 't', 'v', '\\', '"':
		l.advance()
		return true
	default:
		l.errorHandler(l.position.Snapshot(l.offset), `unknown escape: \`+string(l.char))
		return false
	}
}

func (l *lexer) readIdent() string {
	l.readNewLine = true

	start := l.offset
	for isLetter(l.char) || isDecimal(l.char) {
		l.advance()
	}

	if isSpecialChar(l.char) {
		l.advance()
	}

	return string(l.source[start:l.offset])
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
	bytes, err := io.ReadAll(input)
	if err != nil {
		panic("could not read from input" + err.Error())
	}

	l := &lexer{
		source:       bytes,
		errorHandler: errHandler,
		offset:       -1,
		char:         ' ',
		position:     token.NewPosition(),
	}

	return l
}
