package object

func NewEnvironment() *Environment {
	s := make(map[string]Object)
    return &Environment{store: s, outer: nil}
}

func NewEnclosedEnviroment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}

type Environment struct {
	store map[string]Object
    // outer is reference to another environment.
    // when a function is called, the current environment is set as the 
    // outer scope of the environment to create a environment that extends
    // the original outer environment to used during evaluation
    outer *Environment
}

// Get checks the inner scope for a variable with identifier, name
// if it is not found in the inner scope, the outer scope is checked 
// recursively until it is found or the last scope is reached
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
