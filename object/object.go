package object

import (
	"bytes"
	"fmt"
	"monkey/ast"
)

type (
	BuiltinFn  func(args ...Object) Object
	ObjectType string
)

type Object interface {
	Inspect() string
	Type() ObjectType
}

const (
	NULL_OBJ     = "NULL"
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	STRING_OBJ   = "STRING"
	RETURN_OBJ   = "RETURN"
	ERROR_OBJ    = "ERROR"
	FUNCTION_OBJ = "FUNCTION"
	BUILTIN_OBJ  = "BUILTIN"
	ARRAY_OBJ    = "ARRAY"
)

/*** BuiltIn Object ***/

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Inspect() string  { return "bb" }
func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }

/*** Null Object ***/

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return NULL_OBJ }

/*** Integer Object ***/

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

/*** Boolean Object ***/

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }

/*** String Object ***/

type String struct {
	Value string
}

func (s *String) Inspect() string  { return fmt.Sprintf("%v", s.Value) }
func (s *String) Type() ObjectType { return STRING_OBJ }

/*** Return Object ***/

// Wrapper object for a return value
type Return struct {
	Value Object
}

func (r *Return) Inspect() string  { return r.Value.Inspect() }
func (r *Return) Type() ObjectType { return RETURN_OBJ }

/*** Error Object ***/

type Error struct {
	Message string
}

func (err *Error) Inspect() string  { return fmt.Sprintf("ERROR: %s", err.Message) }
func (err *Error) Type() ObjectType { return ERROR_OBJ }

/*** Function Object ***/

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment //Env which the function was declared
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	out.WriteString("fn(")
	for _, param := range f.Parameters {
		out.WriteString(param.String())
	}

	out.WriteString(")")
	out.WriteString(f.Body.String())

	return out.String()
}

/*** Array Object ***/

type Array struct {
	Value []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	for _, val := range a.Value {
		out.WriteString(val.Inspect())
	}

	return out.String()
}
