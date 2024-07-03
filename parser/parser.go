package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/token"
)

// Hashmap to associate a TokenType with a given precedence
var precedence = map[token.TokenType]int{
	token.EQUALITY: EQUALS,
	token.NOTEQUAL: EQUALS,
	token.LT:       LESSGREATER,
	token.GT:       LESSGREATER,
	token.MINUS:    SUM,
	token.PLUS:     SUM,
	token.SLASH:    PRODUCT,
	token.ASTERISK: PRODUCT,
	token.LPAREN:   CALL,
}

// Order of precedence for expression evaluation
const (
	_ int = iota
	LOWEST
	EQUALS
	LESSGREATER
	SUM
	PRODUCT
	PREFIX
	CALL // highest precedence
)

// Pratt Parser
// Tokens are associated with parser functions
// Depending on their position within the expression (prefix/infix) a parsing function will be invoked
type (
	PrefixParseFn func() ast.Expression
	InfixParseFn  func(left ast.Expression) ast.Expression
)

type Parser struct {
	lexer  *lexer.Lexer
	errors []string

	currToken token.Token
	nextToken token.Token

	infixParsers  map[token.TokenType]InfixParseFn
	prefixParsers map[token.TokenType]PrefixParseFn
}

// Create new *Parser struct with given *Lexer
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lexer:         l,
		errors:        []string{},
		infixParsers:  map[token.TokenType]InfixParseFn{},
		prefixParsers: map[token.TokenType]PrefixParseFn{},
	}

	p.registerPrefixParsers()
	p.registerInfixParsers()

	p.advanceTokens()
	p.advanceTokens()

	return p
}

// Construct the AST which represents the given source code
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}

		p.advanceTokens()
	}

	return program
}

// Determine how to parse the current token into an ast.Statement based on the TokenType
func (p *Parser) parseStatement() ast.Statement {
	var stmt ast.Statement

	switch p.currToken.Type {
	case token.LET:
		stmt = p.parseLetStatement()
	case token.RETURN:
		stmt = p.parseReturnStatement()
	default:
		stmt = p.parseExpressionStatement()
	}

	return stmt
}

// Parse the following sete of Statements into an ast.BlockStatement
// Should be called with the curreToken on a LBRACE, and finish with it on the closing RBRACE
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	blkStmt := &ast.BlockStatement{
		Token:      p.currToken,
		Statements: []ast.Statement{},
	}

	p.advanceTokens()

	for !p.currTokenIs(token.RBRACE) && !p.currTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			blkStmt.Statements = append(blkStmt.Statements, stmt)
		}

		p.advanceTokens()
	}

	return blkStmt
}

// Construct the ast.Node to represent a valid LetStatement
func (p *Parser) parseLetStatement() *ast.LetStatement {
	ls := &ast.LetStatement{
		Token: p.currToken, // Let token
	}

	if !p.expectPeek(token.IDENTIFIER) {
		return nil
	}

	ls.Name = &ast.Identifier{
		Token: p.currToken,
		Name:  p.currToken.Literal,
	}

	if !p.expectPeek(token.ASSIGN) {
		return nil
	}

	p.advanceTokens() // Skip over ASSIGN token
	ls.Value = p.parseExpression(LOWEST)

	// Optional semicolon to end statement
	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceTokens()
	}

	return ls
}

// Parse the current token as  ReturnStatement
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	rs := &ast.ReturnStatement{
		Token: p.currToken,
	}

	p.advanceTokens()
	rs.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceTokens()
	}

	return rs
}

// Parse the current token as the start of an ExpressionStatement
// The Expression could be of any ast.Expression type
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	es := &ast.ExpressionStatement{
		Token: p.currToken,
	}

	// LOWEST precedence used since there is nothing to compare yet
	// As expression has not yet begun to be parsed, use the LOWEST precedence to initialize parsing of expression
	es.Expr = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.advanceTokens()
	}

	return es
}