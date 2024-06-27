package lexer

import "monkey/token"

type Lexer struct {
	input   string
	currPos int
	nextPos int
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
	}

	l.advanceChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	l.eatWhitespace()

	// Multi byte tokens
	if isNumber(l.ch) {
		num := l.readWord(isNumber)
		return newTokenStr(token.NUMBER, num)
	} else if isValidChar(l.ch) {
		str := l.readWord(isValidChar)
		tokType := token.GetTokenType(str)
		return newTokenStr(tokType, str)
	}

	// Double byte tokens
	if l.ch == '=' {
		if l.peekCharIs('=') {
			l.advanceChar()
			return newTokenStr(token.EQUALITY, "==")
		}

		l.advanceChar()
		return newTokenStr(token.ASSIGN, "=")
	} else if l.ch == '!' {
		if l.peekCharIs('=') {
			l.advanceChar()
			return newTokenStr(token.NOTEQUAL, "!=")
		}

		l.advanceChar()
		return newTokenStr(token.BANG, "!")
	}

	// Single byte tokens (one char in length)
	var tok token.Token
	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '/':
		tok = newToken(token.SLASH, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case '[':
		tok = newToken(token.LBRACKET, l.ch)
	case ']':
		tok = newToken(token.RBRACKET, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case 0:
		tok = newToken(token.EOF, l.ch)
	default:
		tok = newToken(token.ILLEGAL, l.ch)
	}

	l.advanceChar()
	return tok
}

func (l *Lexer) advanceChar() {
	if l.nextPos >= len(l.input) {
		l.ch = 0
		l.currPos = l.nextPos
		return
	}

	l.ch = l.input[l.nextPos]
	l.currPos = l.nextPos
	l.nextPos++
}

// Check if the next char is certain byte
// Advance lexer character if true
func (l *Lexer) peekCharIs(peek byte) bool {
	if l.nextPos >= len(l.input) {
		return false
	}

	if peek != l.input[l.nextPos] {
		return false
	}

	l.advanceChar()
	return true
}

type PredicateFunc func(byte) bool

func (l *Lexer) readWord(pred PredicateFunc) string {
	start := l.currPos
	for pred(l.ch) {
		l.advanceChar()
	}

	return l.input[start:l.currPos]
}

func (l *Lexer) eatWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\r' || l.ch == '\n' {
		l.advanceChar()
	}
}
