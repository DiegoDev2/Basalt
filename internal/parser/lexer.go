package parser

import (
	"unicode"
)

// TokenType defines the type of a token.
type TokenType string

const (
	IDENT       TokenType = "IDENT"
	STRING      TokenType = "STRING"
	LBRACE      TokenType = "LBRACE"
	RBRACE      TokenType = "RBRACE"
	COLON       TokenType = "COLON"
	AT          TokenType = "AT"
	LPAREN      TokenType = "LPAREN"
	RPAREN      TokenType = "RPAREN"
	HTTP_METHOD TokenType = "HTTP_METHOD"
	PATH        TokenType = "PATH"
	EOF         TokenType = "EOF"
	ILLEGAL     TokenType = "ILLEGAL"
)

// Token represents a single token in the DSL.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
	Column  int
}

// Lexer tokenizes the source code.
type Lexer struct {
	input        string
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           byte // current char under examination
	line         int
	column       int
}

// NewLexer creates a new Lexer instance.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input, line: 1, column: 0}
	l.readChar()
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
	l.column++
}

// NextToken returns the next token from the input.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '{':
		tok = l.newToken(LBRACE, l.ch)
	case '}':
		tok = l.newToken(RBRACE, l.ch)
	case ':':
		tok = l.newToken(COLON, l.ch)
	case '@':
		tok = l.newToken(AT, l.ch)
	case '(':
		tok = l.newToken(LPAREN, l.ch)
	case ')':
		tok = l.newToken(RPAREN, l.ch)
	case '/':
		tok.Type = PATH
		tok.Literal = l.readPath()
		return tok
	case '"':
		tok.Type = STRING
		tok.Literal = l.readString()
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			if isHTTPMethod(tok.Literal) {
				tok.Type = HTTP_METHOD
			} else {
				tok.Type = IDENT
			}
			return tok
		} else {
			tok = l.newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) newToken(tokenType TokenType, ch byte) Token {
	return Token{Type: tokenType, Literal: string(ch), Line: l.line, Column: l.column}
}

func (l *Lexer) readIdentifier() string {
	position := l.position
	for isLetter(l.ch) || unicode.IsDigit(rune(l.ch)) || l.ch == '_' || l.ch == '.' {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

func (l *Lexer) readPath() string {
	position := l.position
	for !unicode.IsSpace(rune(l.ch)) && l.ch != '{' && l.ch != '}' && l.ch != '(' && l.ch != ')' && l.ch != 0 {
		l.readChar()
	}
	return l.input[position:l.position]
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(rune(l.ch)) {
		if l.ch == '\n' {
			l.line++
			l.column = 0
		}
		l.readChar()
	}
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func isHTTPMethod(s string) bool {
	switch s {
	case "GET", "POST", "PUT", "DELETE", "PATCH":
		return true
	}
	return false
}
