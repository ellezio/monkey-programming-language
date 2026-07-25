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
		tok = newToken(token.ASSIGN, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
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

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.ReadChar()
	}
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func isLetter(ch rune) bool {
	if utf8.RuneLen(ch) == 1 {
		return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
	}
	return utf8.ValidRune(ch)
}

func isDigit(ch rune) bool {
	return '0' <= ch && ch <= '9'
}
