package lexer

import "monkey/token"

// Create new token based from provided string
func newTokenStr(tokType token.TokenType, lit string) token.Token {
	tok := token.Token{
		Type:    tokType,
		Literal: lit,
	}

	return tok
}

// Create new token based from provided byte converted to string
func newToken(tokType token.TokenType, lit byte) token.Token {
	tok := token.Token{
		Type:    tokType,
		Literal: string(lit),
	}

	return tok
}

func isNumber(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func notQuote(ch byte) bool {
	return ch != '"'
}

func isIdentifier(ch byte) bool {
	lower := 'a' <= ch && ch <= 'z'
	upper := 'A' <= ch && ch <= 'Z'
	special := ch == '_'

	return lower || upper || special
}
