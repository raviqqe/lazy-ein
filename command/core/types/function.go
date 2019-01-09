package types

// Function is a function type.
type Function struct {
	name      string
	arguments []Type
	result    Type
}

// NewFunction creates a function type.
func NewFunction(as []Type, r Type) Function {
	return Function{"", as, r}
}

// NewNamedFunction creates a function type.
func NewNamedFunction(n string, as []Type, r Type) Function {
	return Function{n, as, r}
}

// Name returns a name.
func (f Function) Name() (string, bool) {
	if f.name == "" {
		return "", false
	}

	return f.name, true
}

// Arguments returns arguments.
func (f Function) Arguments() []Type {
	return f.arguments
}

// Result returns arguments.
func (f Function) Result() Type {
	return f.result
}

func (Function) isType() {}
