package parser

import (
	"fmt"
	"monkey/token"
)

func (p *Parser) currTokenIs(tokType token.TokenType) bool { return p.currToken.Type == tokType }
func (p *Parser) peekTokenIs(tokType token.TokenType) bool { return p.nextToken.Type == tokType }

func (p *Parser) ParserErrors() bool {
	if len(p.errors) > 0 {
		for _, err := range p.errors {
			fmt.Printf("Parser Error: %s\n", err)
		}

		return true
	}

	return false
}

func (p *Parser) advanceTokens() {
	p.currToken = p.nextToken
	p.nextToken = p.lexer.NextToken()
}

// Determines if next token is of the expected type.
// Advances token pointers if true, creates parser error otherwise
func (p *Parser) expectPeek(expToken token.TokenType) bool {
	if !p.peekTokenIs(expToken) {
		p.errors = append(p.errors, fmt.Sprintf("expectPeek found unexpected peek, Got=%s Expected=%s. Line %d, Column %d", p.nextToken.Type, expToken, p.nextToken.Position.Line, p.nextToken.Position.Column))
		return false
	}

	p.advanceTokens()
	return true
}

// Get the precedence associated with the currToken
func (p *Parser) currPrecendence() int {
	if prec, ok := precedence[p.currToken.Type]; ok {
		return prec
	}

	return LOWEST
}

// Get the precedence associated with the nextToken
func (p *Parser) peekPrecendence() int {
	if prec, ok := precedence[p.nextToken.Type]; ok {
		return prec
	}

	return LOWEST
}

// Register tokens which can exist at the beginning of an expression (prefix position) with its associated parsing function.
// prefixParsers advance the currToken to sit on the last token associated with its Node
func (p *Parser) registerPrefixParsers() {
	p.prefixParsers[token.IDENTIFIER] = p.parseIndentifier
	p.prefixParsers[token.NUMBER] = p.parseIntLiteral
	p.prefixParsers[token.STRING] = p.parseStringLiteral
	p.prefixParsers[token.IF] = p.parseConditional

	p.prefixParsers[token.LBRACKET] = p.parseArrayLiteral
	p.prefixParsers[token.LBRACE] = p.parseHashLiteral

	// Prefix operators: Creates a ast.PrefixExpression
	p.prefixParsers[token.BANG] = p.parsePrefixExpression
	p.prefixParsers[token.MINUS] = p.parsePrefixExpression

	// Boolean expressions: ast.BoolLiteral
	p.prefixParsers[token.TRUE] = p.parseBoolLiteral
	p.prefixParsers[token.FALSE] = p.parseBoolLiteral

	// LPAREN prefix position indicates a grouped expression: (<expression>)
	p.prefixParsers[token.LPAREN] = p.parseGroupedExpression

	// Function Literal
	p.prefixParsers[token.FUNCTION] = p.parseFnLiteral
}

// Register Tokens to be used as infix operators with their respective parsing function
// <expression> <infix-token> <expression>
func (p *Parser) registerInfixParsers() {
	// LPAREN infix position indicates call expression: <expression>(<expression args>)
	p.infixParsers[token.LPAREN] = p.parseCallExpression

	p.infixParsers[token.LBRACKET] = p.parseIndexExpression

	p.infixParsers[token.LT] = p.parseInfixExpression
	p.infixParsers[token.GT] = p.parseInfixExpression
	p.infixParsers[token.PLUS] = p.parseInfixExpression
	p.infixParsers[token.MINUS] = p.parseInfixExpression
	p.infixParsers[token.SLASH] = p.parseInfixExpression
	p.infixParsers[token.EQUALITY] = p.parseInfixExpression
	p.infixParsers[token.NOTEQUAL] = p.parseInfixExpression
	p.infixParsers[token.ASTERISK] = p.parseInfixExpression
}
