package parser

import (
	"fmt"
	"monkey/token"
)

func (p *Parser) CheckErrors() {
	for _, err := range p.errors {
		fmt.Println(err)
	}
}

func (p *Parser) advanceTokens() {
	p.currToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

func (p *Parser) currTokenIs(tokType token.TokenType) bool { return p.currToken.Type == tokType }
func (p *Parser) peekTokenIs(tokType token.TokenType) bool { return p.nextToken.Type == tokType }

// Determines if next token is of the expected type.
// Advances token pointers if true, adds error otherwise
func (p *Parser) expectPeek(expToken token.TokenType) bool {
	if !p.peekTokenIs(expToken) {
		p.errors = append(p.errors, fmt.Sprintf("expectPeek found unexpected peek, Got=%s Expected=%s", p.nextToken.Type, expToken))
		return false
	}

	p.advanceTokens()
	return true
}

// Get the precedence associated with the parsers currToken
func (p *Parser) currPrecendence() int {
	if prec, ok := precedence[p.currToken.Type]; ok {
		return prec
	}

	return LOWEST
}

// Get the precedence associated with the parsers nextToken
func (p *Parser) peekPrecendence() int {
	if prec, ok := precedence[p.nextToken.Type]; ok {
		return prec
	}

	return LOWEST
}

// Register the TokenType with its associated parsing function
// Tokens which can exist at the beginning of an expression
// prefixParsers advance the token so it sits on the last token associated with its Node
func (p *Parser) registerPrefixParsers() {
	p.prefixParsers[token.IDENTIFIER] = p.parseIndentifier
	p.prefixParsers[token.NUMBER] = p.parseIntLiteral
	p.prefixParsers[token.IF] = p.parseConditional

	// Prefix operators: Creates a PrefixExpression
	p.prefixParsers[token.BANG] = p.parsePrefixExpression
	p.prefixParsers[token.MINUS] = p.parsePrefixExpression

	// Boolean expressions
	p.prefixParsers[token.TRUE] = p.parseBoolLiteral
	p.prefixParsers[token.FALSE] = p.parseBoolLiteral

	// LPAREN if in prefix position indicates a grouped expression
	//(<expression>)
	p.prefixParsers[token.LPAREN] = p.parseGroupedExpression

	// Function Literal
	p.prefixParsers[token.FUNCTION] = p.parseFnLiteral
}

// Register the Tokens which can be used as an infix operator
// <expression> <infix-token> <expression>
func (p *Parser) registerInfixParsers() {
	// LPAREN infix position indicates call expression
	// <expression>(<expression args>)
	p.infixParsers[token.LPAREN] = p.parseCallExpression

	p.infixParsers[token.LT] = p.parseInfixExpression
	p.infixParsers[token.GT] = p.parseInfixExpression
	p.infixParsers[token.PLUS] = p.parseInfixExpression
	p.infixParsers[token.MINUS] = p.parseInfixExpression
	p.infixParsers[token.SLASH] = p.parseInfixExpression
	p.infixParsers[token.EQUALITY] = p.parseInfixExpression
	p.infixParsers[token.NOTEQUAL] = p.parseInfixExpression
	p.infixParsers[token.ASTERISK] = p.parseInfixExpression
}
