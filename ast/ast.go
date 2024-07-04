package ast

import "bytes"

type Node interface {
	String() string
	TokenLiteral() string
}

// Interface to distinguish the Expression/Statement types
// Empty func only used for LSP help
type Statement interface {
	Node
	statment()
}

type Expression interface {
	Node
	expression()
}

// List of Statements which represents the source codes structure
// Each statement must properly nest any Expressions for capture the semantics of the source code
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}

	return p.Statements[0].TokenLiteral()
}

func (p *Program) String() string {
	var out bytes.Buffer

	for _, stmt := range p.Statements {
		out.WriteString(stmt.String())
		out.WriteString("\n")
		out.WriteString("\n")
	}

	return out.String()
}
