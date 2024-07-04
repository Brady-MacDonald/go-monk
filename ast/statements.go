package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

/*** Block Statement ***/
type BlockStatement struct {
	Token      token.Token // Opening LBRACE of statement {
	Statements []Statement
}

func (bs *BlockStatement) statment()            {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{")

	if len(bs.Statements) > 0 {
		for _, stmt := range bs.Statements {
			out.WriteString(stmt.String())
		}
	}

	out.WriteString("}")
	return out.String()
}

/*** Let Statement ***/
type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statment() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	letStr := fmt.Sprintf("let %s = %s;", ls.Name.String(), ls.Value.String())
	return letStr
}

/*** Return Statement ***/
type ReturnStatement struct {
	Token     token.Token
	ReturnExp Expression // Expression to be returned
}

func (rs *ReturnStatement) statment() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	returnStr := fmt.Sprintf("return %s", rs.ReturnExp.String())
	return returnStr
}

/*** Expression Statement ***/
type ExpressionStatement struct {
	Token token.Token
	Expr  Expression
}

func (es *ExpressionStatement) statment() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	return es.Expr.String()
}
