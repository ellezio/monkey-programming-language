package lexer

import (
	"monkey/token"
	"unicode/utf8"
)

type Lexer struct {
	input        string
	position     int  // current position in input
	readPosition int  // next position - serves purpose of peek
	ch           rune // current character on which position points
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.ReadChar() // It's called in order to return Lexer in working state - pointing to the first character
	return l
}

func (l *Lexer) ReadChar() {
	rb := 1
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch, rb = utf8.DecodeRuneInString(l.input[l.readPosition:])
	}

	l.position = l.readPosition
	l.readPosition += rb
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if tok = l.makeTwoCharToken(token.EQ, '='); tok == token.Nil {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if tok = l.makeTwoCharToken(token.NOT_EQ, '='); tok == token.Nil {
			tok = newToken(token.BANG, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.ReadChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.ReadChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.ReadChar()
		if l.ch == rune('"') || l.ch == 0 {
			break
		}
	}

	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadChar()
	}
}

func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		ch, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
		return ch
	}
}

func (l *Lexer) makeTwoCharToken(tokenType token.TokenType, nextCh rune) token.Token {
	tok := token.Nil
	if l.peekChar() == nextCh {
		ch := l.ch
		l.ReadChar()
		tok.Type = tokenType
		tok.Literal = string(ch) + string(l.ch)
	}

	return tok
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
	if ch < utf8.RuneSelf {
		return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
	}
	return ch != utf8.RuneError
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
