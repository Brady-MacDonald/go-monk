package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

// List of TokenTypes
const (
	EQUALITY  = "EQUALITY"    // ==
	NOTEQUAL  = "NOTEQUAL"    // !=
	LT        = "LESSTHAN"    // <
	GT        = "GREATERTHAN" // >
	ASSIGN    = "ASSIGN"      // =
	BANG      = "BANG"        // !
	PLUS      = "PLUS"        // +
	MINUS     = "MINUS"       // -
	SLASH     = "SLASH"       // /
	ASTERISK  = "ASTERISK"    // *
	LBRACE    = "LBRACE"      // (
	RBRACE    = "RBRACE"      // )
	LPAREN    = "LPAREN"      // {
	RPAREN    = "RPAREN"      // }
	LBRACKET  = "LBRACKET"    // [
	RBRACKET  = "RBRACKET"    // ]
	SEMICOLON = "SEMICOLON"   // ;

	IF       = "IF"
	ELSE     = "ELSE"
	RETURN   = "RETURN"
	LET      = "LET"
	FUNCTION = "FUNCTION"
	TRUE     = "TRUE"
	FALSE    = "FALSE"

	NUMBER     = "NUMBER"
	IDENTIFIER = "IDENTIFIER"

	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"
)

var keywords = map[string]TokenType{
	"if":     IF,
	"let":    LET,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
	"fn":     FUNCTION,
}

// Determine if the given identifier is a language defined keyword
func GetTokenType(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
