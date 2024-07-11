package object

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"monkey/ast"
)

type (
	BuiltinFn  func(args ...Object) Object
	ObjectType string
)

type Hashable interface {
	HashKey() HashKey
}

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
	HASH_OBJ     = "HASH"
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

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	hashKey := HashKey{
		Type: INTEGER_OBJ,
		Key:  uint64(i.Value),
	}

	return hashKey
}

/*** Boolean Object ***/

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%v", b.Value) }
func (b *Boolean) HashKey() HashKey {
	hashKey := HashKey{
		Type: BOOLEAN_OBJ,
	}

	if b.Value {
		hashKey.Key = 1
	} else {
		hashKey.Key = 0
	}

	return hashKey
}

/*** String Object ***/

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf("\"%v\"", s.Value) }
func (s *String) HashKey() HashKey {
	hashKey := HashKey{
		Type: STRING_OBJ,
	}
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	hashKey.Key = h.Sum64()

	return hashKey
}

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

	out.WriteString("[")
	for _, val := range a.Value {
		out.WriteString(val.Inspect())
		out.WriteString(",")
	}
	out.WriteString("]")

	return out.String()
}

/*** Hash ***/

type HashKey struct {
	Type string
	Key  uint64
}

// Used to track the original key, not the hashed key
// Useful for printing the Hashtable
type HashPair struct {
	Key Object
	Val Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer

	out.WriteString("{")
	for _, pair := range h.Pairs {
		out.WriteString(fmt.Sprintf("%s: %s,", pair.Key.Inspect(), pair.Val.Inspect()))
	}
	out.WriteString("}")

	return out.String()
}
