package token

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Multi byte tokens

	NUMBER     = "Number"
	STRING     = "String"
	IDENTIFIER = "Identifier"

	// Two byte tokens

	EQUALITY = "Equality" // ==
	NOTEQUAL = "NotEqual" // !=

	// Single byte tokens

	LT        = "LessThan"      // <
	GT        = "GreaterThan"   // >
	ASSIGN    = "Assign"        // =
	BANG      = "Bang"          // !
	PLUS      = "Plus"          // +
	MINUS     = "Minus"         // -
	SLASH     = "Slash"         // /
	ASTERISK  = "Asterisk"      // *
	LBRACE    = "Left-Brace"    // {
	RBRACE    = "Right-Brace"   // }
	LPAREN    = "Left-Paren"    // (
	RPAREN    = "Right-Paren"   // )
	LBRACKET  = "Left-Bracket"  // [
	RBRACKET  = "Right-Bracket" // ]
	SEMICOLON = "Semicolon"     // ;
	COLON     = "Colon"         // :
	COMMA     = "Comma"         // ,

	// Keywords

	IF       = "If"
	ELSE     = "Else"
	RETURN   = "Return"
	LET      = "Let"
	FUNCTION = "Function"
	TRUE     = "True"
	FALSE    = "False"

	ILLEGAL = "Illegal"
	EOF     = "EOF"
)

// Map the language keywords, to their corresponding TokenType
var keywords = map[string]TokenType{
	"if":     IF,
	"let":    LET,
	"else":   ELSE,
	"true":   TRUE,
	"false":  FALSE,
	"return": RETURN,
	"fn":     FUNCTION,
}

// Determine if the given identifier is a language keyword
func GetTokenType(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}

	return IDENTIFIER
}
