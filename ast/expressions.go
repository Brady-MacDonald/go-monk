package ast

import (
	"bytes"
	"fmt"
	"monkey/token"
)

/*** Integer Literal ***/

type IntLiteral struct {
	Token token.Token
	Value int64
}

func (i *IntLiteral) expression()          {}
func (i *IntLiteral) String() string       { return fmt.Sprintf("%d", i.Value) }
func (i *IntLiteral) TokenLiteral() string { return i.Token.Literal }

/*** Boolean Literal ***/

type BoolLiteral struct {
	Token token.Token
	Value bool
}

func (b *BoolLiteral) expression()          {}
func (b *BoolLiteral) String() string       { return fmt.Sprintf("%v", b.Value) }
func (b *BoolLiteral) TokenLiteral() string { return b.Token.Literal }

/*** String Literal ***/

type StringLiteral struct {
	Token token.Token
	Value string
}

func (s *StringLiteral) expression()          {}
func (s *StringLiteral) String() string       { return s.Value }
func (s *StringLiteral) TokenLiteral() string { return s.Token.Literal }

/*** Array Literal ***/

type ArrayLiteral struct {
	Token    token.Token // [
	Elements []Expression
}

func (a *ArrayLiteral) expression()          {}
func (a *ArrayLiteral) TokenLiteral() string { return a.Token.Literal }
func (a *ArrayLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("[")
	for _, elem := range a.Elements {
		out.WriteString(elem.String())
		out.WriteString(",")
	}
	out.WriteString("]")

	return out.String()
}

/*** Hash Literal ***/

type HashLiteral struct {
	Token token.Token // {
	Pairs map[Expression]Expression
}

func (h *HashLiteral) expression()          {}
func (h *HashLiteral) TokenLiteral() string { return h.Token.Literal }
func (h *HashLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("{")
	for key, val := range h.Pairs {
		out.WriteString(fmt.Sprintf("   %s: %s,\n", key.String(), val.String()))
	}
	out.WriteString("}")

	return out.String()
}

/*** Identifier ***/

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expression()          {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

/*** Index Expression ***/

type IndexExpression struct {
	Token token.Token //[
	Left  Expression
	Index Expression
}

func (i *IndexExpression) expression()          {}
func (i *IndexExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IndexExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(i.Left.String())
	out.WriteString("[")
	out.WriteString(i.Index.String())
	out.WriteString("])")

	return out.String()
}

/*** Prefix Expression ***/

type PrefixExpression struct {
	Token    token.Token
	Operator string // ! or -
	Operand  Expression
}

func (pe *PrefixExpression) expression() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}

func (pe *PrefixExpression) String() string {
	var out bytes.Buffer

	// Brackets around Operator/Right to indicate the grouping of operator/operand
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Operand.String())
	out.WriteString(")")

	return out.String()
}

/*** Infix Expression ***/

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expression()          {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()

}

/*** If Expression ***/

// if keyword is an expression, produces a value like a ternary
type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement // Executed for 'if' block
	Alternative *BlockStatement // Executed as 'else' block
}

func (i *IfExpression) expression()          {}
func (i *IfExpression) TokenLiteral() string { return i.Token.Literal }
func (i *IfExpression) String() string {
	var out bytes.Buffer

	out.WriteString("if(")
	out.WriteString(i.Condition.String())
	out.WriteString(")")
	out.WriteString(i.Consequence.String())

	// else block is optional
	if i.Alternative != nil {
		out.WriteString("else")
		out.WriteString(i.Alternative.String())
	}

	return out.String()
}

/*** Function Literal ***/

type FnLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (fl *FnLiteral) expression()          {}
func (fl *FnLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FnLiteral) String() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for idx, ident := range fl.Parameters {
		out.WriteString(ident.String())
		if idx < len(fl.Parameters)-1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")
	out.WriteString(fl.Body.String())

	return out.String()
}

/*** Call Expression ***/

type CallExpression struct {
	Token token.Token // ( token
	Fn    Expression  // Identifier or FnLiteral
	Args  []Expression
}

func (ce *CallExpression) expression()          {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer

	out.WriteString("fn")
	out.WriteString("(")
	for _, arg := range ce.Args {
		out.WriteString(arg.String())
		out.WriteString(",")
	}
	out.WriteString(")")

	return out.String()
}
