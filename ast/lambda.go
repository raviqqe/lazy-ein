package ast

import "github.com/raviqqe/stg/types"

// Lambda is a lambda form.
type Lambda struct {
	arguments     []Argument
	body          Expression
	resultType    types.Type
	freeVariables []Argument
	updatable     bool
}

// NewLambda creates a new lambda form.
func NewLambda(vs []Argument, u bool, as []Argument, e Expression, t types.Type) Lambda {
	return Lambda{as, e, t, vs, u}
}

// ArgumentNames returns argument names.
func (l Lambda) ArgumentNames() []string {
	return argumentsToNames(l.arguments)
}

// ArgumentTypes returns argument types.
func (l Lambda) ArgumentTypes() []types.Type {
	return argumentsToTypes(l.arguments)
}

// Body returns a body expression.
func (l Lambda) Body() Expression {
	return l.body
}

// ResultType returns a result type.
func (l Lambda) ResultType() types.Type {
	return l.resultType
}

// FreeVariableNames returns free varriable names.
func (l Lambda) FreeVariableNames() []string {
	return argumentsToNames(l.freeVariables)
}

// FreeVariableTypes returns free varriable types.
func (l Lambda) FreeVariableTypes() []types.Type {
	return argumentsToTypes(l.freeVariables)
}

// IsUpdatable returns true if the lambda form is updatable, or false otherwise.
func (l Lambda) IsUpdatable() bool {
	return l.updatable
}

func argumentsToNames(as []Argument) []string {
	ss := make([]string, 0, len(as))

	for _, a := range as {
		ss = append(ss, a.Name())
	}

	return ss
}

func argumentsToTypes(as []Argument) []types.Type {
	ts := make([]types.Type, 0, len(as))

	for _, a := range as {
		ts = append(ts, a.Type())
	}

	return ts
}
