package ast

// ConstructorApplication is a constructor application.
type ConstructorApplication struct {
	name      string
	arguments []Atom
}

// NewConstructorApplication creates a constructor application.
func NewConstructorApplication(n string, as []Atom) ConstructorApplication {
	return ConstructorApplication{n, as}
}

// Name returns a name.
func (c ConstructorApplication) Name() string {
	return c.name
}

// Arguments returns arguments.
func (c ConstructorApplication) Arguments() []Atom {
	return c.arguments
}

func (ConstructorApplication) isExpression() {}
