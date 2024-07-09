package parser

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"strconv"
)

// Begin the parsing of an Expression, initiated using the LOWEST precedence
func (p *Parser) parseExpression(precedence int) ast.Expression {
	// Get the parser function associated to the current token
	prefixFn := p.prefixParsers[p.currToken.Type]
	if prefixFn == nil {
		// Do not know how to begin parsing an expression with this TokenType
		p.errors = append(p.errors, fmt.Sprintf("No prefix parser function found for token %s", p.currToken.Type))
		return nil
	}

	expr := prefixFn()

	// Iterate while the next token is not a SEMICOLON
	// And the nextToken has a higher precedence than the current one
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecendence() {
		// Locate infix parsing func for next token
		infixFn := p.infixParsers[p.nextToken.Type]
		if infixFn == nil {
			// If there's no infix parseFn then we are done parsing expression
			// nextToken maybe a let/return token
			// maybe an int/boolen (some expression which is terminal and not recursive)
			return expr
		}

		p.advanceTokens() //advance tokens to currToken sits on the infix operator
		expr = infixFn(expr)
	}

	return expr
}

// Parse the current token as a PrefixExpression.
// <operator><expression>
func (p *Parser) parsePrefixExpression() ast.Expression {
	prefix := &ast.PrefixExpression{
		Token:    p.currToken,
		Operator: p.currToken.Literal, // ! or -
	}

	p.advanceTokens()

	prefix.Operand = p.parseExpression(PREFIX)
	return prefix
}

// Parse the current token as the operator to an InfixExpression.
// Left side of infix expression already parsed.
// <expression><operator><expression>
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	infix := &ast.InfixExpression{
		Left:     left,
		Token:    p.currToken,
		Operator: p.currToken.Literal,
	}

	prec := p.currPrecendence()
	p.advanceTokens()

	infix.Right = p.parseExpression(prec)
	return infix
}

// Parse the current token as an Indentifier
func (p *Parser) parseIndentifier() ast.Expression {
	ident := &ast.Identifier{
		Token: p.currToken,
		Value: p.currToken.Literal,
	}

	return ident
}

// Parse the current token as an IntLiteral
func (p *Parser) parseIntLiteral() ast.Expression {
	intLiteral := &ast.IntLiteral{
		Token: p.currToken,
	}

	val, err := strconv.ParseInt(p.currToken.Literal, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Unable to parse %s to int literal", p.currToken.Literal))
		return nil
	}

	intLiteral.Value = val
	return intLiteral
}

// Parse the current token as an BoolLiteral
func (p *Parser) parseBoolLiteral() ast.Expression {
	boolExp := &ast.BoolLiteral{
		Token: p.currToken,
	}

	val, err := strconv.ParseBool(p.currToken.Literal)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("Unable to parse %s as a Boolean", p.currToken.Literal))
		return nil
	}

	boolExp.Value = val
	return boolExp
}

// Current token is a LPAREN in the prefix position
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.advanceTokens()
	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// Parse currToken as in IfExpression
func (p *Parser) parseConditional() ast.Expression {
	ifExpr := &ast.IfExpression{
		Token: p.currToken,
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.advanceTokens()
	ifExpr.Condition = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	ifExpr.Consequence = p.parseBlockStatement()

	// Optional else clause
	if p.peekTokenIs(token.ELSE) {
		p.advanceTokens()
		if !p.expectPeek(token.LBRACE) {
			return nil
		}

		ifExpr.Alternative = p.parseBlockStatement()
	}

	return ifExpr
}

// Parse currToken as a FnLiteral
func (p *Parser) parseFnLiteral() ast.Expression {
	fn := &ast.FnLiteral{
		Token:      p.currToken,
		Parameters: []*ast.Identifier{},
	}

	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	fn.Parameters = p.parseFnParams()

	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	fn.Body = p.parseBlockStatement()

	return fn
}

// Helper for parsing a functions params
func (p *Parser) parseFnParams() []*ast.Identifier {
	params := []*ast.Identifier{}

	// No params
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return params
	}

	// First param
	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	ident := p.parseIndentifier().(*ast.Identifier)
	params = append(params, ident)

	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		if !p.expectPeek(token.IDENTIFIER) {
			return nil
		}

		ident := p.parseIndentifier().(*ast.Identifier)
		params = append(params, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return params
}

// Current token is LPAREN in an infix position
// Previous should be an Identifier/FnLiteral
func (p *Parser) parseCallExpression(fn ast.Expression) ast.Expression {
	callExpr := &ast.CallExpression{
		Token: p.currToken,
		Fn:    fn,
	}

	callExpr.Args = p.parseCallArgs()
	return callExpr
}

// Helper to parse the expressions passed as args to fn invocation
func (p *Parser) parseCallArgs() []ast.Expression {
	args := []ast.Expression{}

	// Empty args
	if p.peekTokenIs(token.RPAREN) {
		p.advanceTokens()
		return args
	}

	p.advanceTokens()
	firstArg := p.parseExpression(LOWEST)
	args = append(args, firstArg)

	for p.peekTokenIs(token.COMMA) {
		p.advanceTokens()
		p.advanceTokens()
		argExpr := p.parseExpression(LOWEST)
		args = append(args, argExpr)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}
