package ast

import "github.com/raviqqe/stg/types"

// Lambda is a lambda form.
type Lambda struct {
	arguments     []Argument
	body          Expression
	resultType    types.Type
	freeVariables []string
	updatable     bool
}

// NewLambda creates a new lambda form.
func NewLambda(fs []string, u bool, as []Argument, e Expression, t types.Type) Lambda {
	return Lambda{as, e, t, fs, u}
}

// ArgumentNames returns argument names.
func (l Lambda) ArgumentNames() []string {
	ss := make([]string, 0, len(l.arguments))

	for _, a := range l.arguments {
		ss = append(ss, a.Name())
	}

	return ss
}

// ArgumentTypes returns argument types.
func (l Lambda) ArgumentTypes() []types.Type {
	ts := make([]types.Type, 0, len(l.arguments))

	for _, a := range l.arguments {
		ts = append(ts, a.Type())
	}

	return ts
}

// Body returns a body expression.
func (l Lambda) Body() Expression {
	return l.body
}

// ResultType returns a result type.
func (l Lambda) ResultType() types.Type {
	return l.resultType
}

// FreeVariables returns a body expression.
func (l Lambda) FreeVariables() []string {
	return l.freeVariables
}

// IsUpdatable returns true if the lambda form is updatable, or false otherwise.
func (l Lambda) IsUpdatable() bool {
	return l.updatable
}
