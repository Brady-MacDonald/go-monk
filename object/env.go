package object

/*** Environment ***/

// Environment used to map identifiers to values
//
// Allows to reference outer Environments creating closures and variable scoping
type Environment struct {
	store map[string]Object
	outer *Environment
}

// Create a new empty env
func NewEnvironment() *Environment {
	return &Environment{
		outer: nil,
		store: map[string]Object{},
	}
}

// Create new Environment and wrap the enclosing scope
func NewEnclosingEnvironment(outer *Environment) *Environment {
	return &Environment{
		outer: outer,
		store: map[string]Object{},
	}
}

// Bind an identifier to the given value
func (e *Environment) Set(key string, val Object) {
	e.store[key] = val
}

// Checks the current Environment for a given identifier
// If not found the outer Environment(s) is checked
func (e *Environment) Get(key string) Object {
	val, ok := e.store[key]
	if !ok && e.outer != nil {
		return e.outer.Get(key)
	}

	return val
}
