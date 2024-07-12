package lexer

import "monkey/token"

type Lexer struct {
	input   string
	line    int
	column  int
	currPos int // Index of current char in input string
	nextPos int // Index of next char to examine
	ch      byte
}

func New(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
	}

	l.advanceChar()
	return l
}

// Get the next token from the source code
func (l *Lexer) NextToken() token.Token {
	l.eatWhitespace()

	pos := token.Position{
		Filename: "test.monk",
		Line:     l.line,
		Column:   l.column,
	}

	// Multi byte token if first character is a number/letter/quote
	if isNumber(l.ch) {
		num := l.readWord(isNumber)
		return newTokenStr(token.NUMBER, num, pos)

	} else if isIdentifier(l.ch) {
		str := l.readWord(isIdentifier)
		tokType := token.GetTokenType(str)
		return newTokenStr(tokType, str, pos)

	} else if !notQuote(l.ch) { //not(notQuote) == isQuote :/
		l.advanceChar()
		str := l.readWord(notQuote)
		l.advanceChar()
		return newTokenStr(token.STRING, str, pos)
	}

	// Double byte tokens
	if l.ch == '=' {
		if l.peekCharIs('=') {
			l.advanceChar()
			return newTokenStr(token.EQUALITY, "==", pos)
		}

		l.advanceChar()
		return newToken(token.ASSIGN, '=', pos)
	} else if l.ch == '!' {
		if l.peekCharIs('=') {
			l.advanceChar()
			return newTokenStr(token.NOTEQUAL, "!=", pos)
		}

		l.advanceChar()
		return newToken(token.BANG, '!', pos)
	}

	// Single byte tokens (one char in length)
	var tok token.Token
	switch l.ch {
	case '+':
		tok = newToken(token.PLUS, l.ch, pos)
	case '-':
		tok = newToken(token.MINUS, l.ch, pos)
	case '*':
		tok = newToken(token.ASTERISK, l.ch, pos)
	case '/':
		tok = newToken(token.SLASH, l.ch, pos)
	case '<':
		tok = newToken(token.LT, l.ch, pos)
	case '>':
		tok = newToken(token.GT, l.ch, pos)
	case '(':
		tok = newToken(token.LPAREN, l.ch, pos)
	case ')':
		tok = newToken(token.RPAREN, l.ch, pos)
	case '{':
		tok = newToken(token.LBRACE, l.ch, pos)
	case '}':
		tok = newToken(token.RBRACE, l.ch, pos)
	case '[':
		tok = newToken(token.LBRACKET, l.ch, pos)
	case ']':
		tok = newToken(token.RBRACKET, l.ch, pos)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch, pos)
	case ':':
		tok = newToken(token.COLON, l.ch, pos)
	case ',':
		tok = newToken(token.COMMA, l.ch, pos)
	case 0:
		tok = newToken(token.EOF, l.ch, pos)
	default:
		tok = newToken(token.ILLEGAL, l.ch, pos)
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
	if l.ch == '\n' {
		l.line++
		l.column = -1
	}

	l.currPos = l.nextPos
	l.nextPos++
	l.column++
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

// Use the given predicate to read up until the end of the number/identifier
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
