package ast

import "github.com/raviqqe/lazy-ein/command/types"

// Variable is a variable.
type Variable struct {
	name string
}

// NewVariable creates a variable.
func NewVariable(s string) Variable {
	return Variable{s}
}

// Name returns a name.
func (v Variable) Name() string {
	return v.name
}

// ConvertExpressions converts expressions.
func (v Variable) ConvertExpressions(f func(Expression) Expression) Expression {
	return f(v)
}

// VisitTypes visits types.
func (v Variable) VisitTypes(f func(types.Type) error) error {
	return nil
}

func (Variable) isExpression() {}
