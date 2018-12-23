package ast

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

// ConvertExpressions visits expressions.
func (v Variable) ConvertExpressions(f func(Expression) Expression) node {
	return f(v)
}

func (Variable) isExpression() {}
