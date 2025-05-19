package object

type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment creates a new Environment with an empty store and no outer environment.
// This is typically used to create the global environment for a program.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment creates a new Environment with an empty store and the given outer environment.
// This is used to create a new scope (e.g., for function calls) that has access to variables
// in the outer scope through the outer environment.
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
