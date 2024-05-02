package lexer

import (
	"strings"

	"github.com/dominicgaliano/interpreter-demo/token"
)

type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // initialize Lexer state
	return l
}

func (l *Lexer) readChar() {
	// set ch to ASCII NUL on end of file
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition += 1
}

func (l *Lexer) readIdentifier() string {
	var builder strings.Builder

	for isLetter(l.ch) {
		builder.WriteByte(l.ch)
		l.readChar()
	}

	return builder.String()
}

func (l *Lexer) readNumber() string {
	var builder strings.Builder

	for isDigit(l.ch) {
		builder.WriteByte(l.ch)
		l.readChar()
	}

	return builder.String()
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		tok = newToken(token.ASSIGN, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		tok = newToken(token.BANG, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case 0:
		tok = newToken(token.EOF, 0)
	default:
		// token is not a special character,
		if isLetter(l.ch) {
			// parse identifier
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdentifier(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			// parse integer literal
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isWhitespace(ch byte) bool {
	asciiWhitespaceMap := map[byte]bool{
		32: true, // Space
		9:  true, // Horizontal tab
		10: true, // Newline
		11: true, // Vertical tab
		12: true, // Form feed
		13: true, // Carriage return
	}

	return asciiWhitespaceMap[ch]
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	if ch == 0 {
		return token.Token{Type: tokenType, Literal: ""}
	}
	return token.Token{Type: tokenType, Literal: string(ch)}
}

func (l *Lexer) skipWhitespace() {
	for isWhitespace(l.ch) {
		l.readChar()
	}
}
