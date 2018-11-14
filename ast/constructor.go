package ast

// Constructor is a constructor application.
type Constructor struct {
	name      string
	arguments []Atom
}

// NewConstructor creates a constructor application.
func NewConstructor(n string, as []Atom) Constructor {
	return Constructor{n, as}
}

// Name returns a name.
func (c Constructor) Name() string {
	return c.name
}

// Arguments returns arguments.
func (c Constructor) Arguments() []Atom {
	return c.arguments
}

func (Constructor) isExpression() {}
