package ast

// Variable is a variable.
type Variable string

// NewVariable creates a variable.
func NewVariable(s string) Variable {
	return Variable(s)
}

// Name returns a name.
func (v Variable) Name() string {
	return string(v)
}

func (Variable) isAtom()       {}
func (Variable) isExpression() {}
