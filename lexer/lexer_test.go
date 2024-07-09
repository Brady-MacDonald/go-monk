package lexer

import (
	"monkey/token"
	"testing"
)

func TestNextToken(t *testing.T) {
	input := `
		let five = 5;
		let ten = 10;

		let add = fn(x, y) {
			x + y;
		};

		let result = add(five, ten);
		!-/*5;
		5 < 10 > 5;

		if (5 < 10) {
			return true;
		} else {
			return false;
		}
		10 == 10;
		10 != 9;
		[1, 2];
	`

	tests := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{token.LET, "let"},
		{token.IDENTIFIER, "five"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "ten"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "x"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENTIFIER, "x"},
		{token.PLUS, "+"},
		{token.IDENTIFIER, "y"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.LET, "let"},
		{token.IDENTIFIER, "result"},
		{token.ASSIGN, "="},
		{token.IDENTIFIER, "add"},
		{token.LPAREN, "("},
		{token.IDENTIFIER, "five"},
		{token.COMMA, ","},
		{token.IDENTIFIER, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.NUMBER, "5"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.NUMBER, "10"},
		{token.EQUALITY, "=="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.NUMBER, "10"},
		{token.NOTEQUAL, "!="},
		{token.NUMBER, "9"},
		{token.SEMICOLON, ";"},
		// {token.STRING, "foobar"},
		// {token.STRING, "foo bar"},
		{token.LBRACKET, "["},
		{token.NUMBER, "1"},
		{token.COMMA, ","},
		{token.NUMBER, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLON, ";"},
		// {token.LBRACE, "{"},
		// {token.STRING, "foo"},
		// {token.COLON, ":"},
		// {token.STRING, "bar"},
		// {token.RBRACE, "}"},
		// {token.EOF, ""},
	}

	l := New(input)
	for idx, tt := range tests {
		tok := l.NextToken()
		if tok.Type != tt.expType {
			t.Fatalf("test[%d]: Invalid TokenType. Got %s, Expected %s", idx, tok.Type, tt.expType)
		}
		if tok.Literal != tt.expLiteral {
			t.Fatalf("test[%d]: Invalid Literal. Got %s, Expected %s", idx, tok.Literal, tt.expLiteral)
		}
	}
}

func TestLexer(t *testing.T) {
	input := "let x = 5"

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

func TestLexerMath(t *testing.T) {
	input := "+-*/"
	expected := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{
			expType:    token.PLUS,
			expLiteral: "+",
		},
		{
			expType:    token.MINUS,
			expLiteral: "-",
		},
		{
			expType:    token.ASTERISK,
			expLiteral: "*",
		},
		{
			expType:    token.SLASH,
			expLiteral: "/",
		},
	}

	l := New(input)
	for idx, tt := range expected {
		token := l.NextToken()

		if token.Type != tt.expType {
			t.Fatalf("test[%d]: Incorrect TokenType. Got %s, expected %s", idx, token.Type, tt.expType)
		}
		if token.Literal != tt.expLiteral {
			t.Fatalf("test[%d]: Incorrect Literal. Got %s, expected %s", idx, token.Literal, tt.expLiteral)
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

func TestMultiByteTokens(t *testing.T) {
	input := `
        tester
        123
        name
        ==
        age
        !=
    `

	expected := []struct {
		expType    token.TokenType
		expLiteral string
	}{
		{
			expType:    token.IDENTIFIER,
			expLiteral: "tester",
		},
		{
			expType:    token.NUMBER,
			expLiteral: "123",
		},
		{
			expType:    token.IDENTIFIER,
			expLiteral: "name",
		},
		{
			expType:    token.EQUALITY,
			expLiteral: "==",
		},
		{
			expType:    token.IDENTIFIER,
			expLiteral: "age",
		},
		{
			expType:    token.NOTEQUAL,
			expLiteral: "!=",
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
