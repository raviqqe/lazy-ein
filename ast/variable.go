package ast

// Variable is a variable.
type Variable string

// NewVariable creates a new variable.
func NewVariable(s string) Variable {
	return Variable(s)
}
