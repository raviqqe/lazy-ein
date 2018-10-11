package ast

// Lambda is a lambda form.
type Lambda struct {
	arguments     []string
	body          Expression
	freeVariables []string
	updatable     bool
}

// NewLambda creates a new lambda form.
func NewLambda(fs []string, u bool, as []string, e Expression) Lambda {
	return Lambda{as, e, fs, u}
}

// Arguments returns arguments.
func (l Lambda) Arguments() []string {
	return l.arguments
}

// Body returns a body expression.
func (l Lambda) Body() Expression {
	return l.body
}

// FreeVariables returns a body expression.
func (l Lambda) FreeVariables() []string {
	return l.freeVariables
}

// IsUpdatable returns true if the lambda form is updatable, or false otherwise.
func (l Lambda) IsUpdatable() bool {
	return l.updatable
}
