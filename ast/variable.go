package ast

// Variable is a variable.
type Variable string

// NewVariable creates a new variable.
func NewVariable(s string) Variable {
	return Variable(s)
}

// Name returns a name.
func (v Variable) Name() string {
	return string(v)
}

func (v Variable) isAtom() {}
