package lexer

import (
	"monkey/token"
	"testing"
)

func TestLexer(t *testing.T) {
	input := `
        let x = 5;
    `
	expected := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{
			expType:    token.LET,
			expLiteral: "let",
		},
		{
			expType:    token.IDENTIFIER,
			expLiteral: "x",
		},
		{
			expType:    token.ASSIGN,
			expLiteral: "=",
		},
		{
			expType:    token.NUMBER,
			expLiteral: "5",
		},
		{
			expType:    token.SEMICOLON,
			expLiteral: ";",
		},
	}

	l := New(input)
	for idx, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("test[%d]: Invalid TokenType. Got %s, Expected %s", idx, tok.Type, tt.expType)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("test[%d]: Invalid Literal. Got %s, Expected %s", idx, tok.Literal, tt.expLiteral)
		}
	}
}

func TestWierd(t *testing.T) {
	input := "let = =="
	expected := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{
			expType:    token.LET,
			expLiteral: "let",
		},
		{
			expType:    token.ASSIGN,
			expLiteral: "=",
		},
		{
			expType:    token.EQUALITY,
			expLiteral: "==",
		},
	}

	l := New(input)
	for idx, tt := range expected {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("test[%d]: Invalid TokenType. Got %s, Expected %s", idx, tok.Type, tt.expType)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("test[%d]: Invalid Literal. Got %s, Expected %s", idx, tok.Literal, tt.expLiteral)
		}
	}
}

func TestLexerSymbols(t *testing.T) {
	input := "(){}[]"

	expected := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{
			expType:    token.LPAREN,
			expLiteral: "(",
		},
		{
			expType:    token.RPAREN,
			expLiteral: ")",
		},
		{
			expType:    token.LBRACE,
			expLiteral: "{",
		},
		{
			expType:    token.RBRACE,
			expLiteral: "}",
		},
		{
			expType:    token.LBRACKET,
			expLiteral: "[",
		},
		{
			expType:    token.RBRACKET,
			expLiteral: "]",
		},
	}

	l := New(input)
	for idx, tt := range expected {
		tok := l.NextToken()

		if tok.Type != tt.expType {
			t.Fatalf("test[%d]: Invalid TokenType. Got %s, Expected %s", idx, tok.Type, tt.expType)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("test[%d]: Invalid Literal. Got %s, Expected %s", idx, tok.Literal, tt.expLiteral)
		}
	}
}
