package ast

import (
	"github.com/ein-lang/ein/command/core/types"
)

// Lambda is a lambda form.
type Lambda struct {
	LambdaDeclaration
	freeVariables []string
	arguments     []string
	body          Expression
}

// NewVariableLambda creates a lambda form.
func NewVariableLambda(vs []Argument, e Expression, t types.Bindable) Lambda {
	return newLambda(vs, nil, e, t)
}

// NewFunctionLambda creates a lambda form.
func NewFunctionLambda(vs []Argument, as []Argument, e Expression, t types.Type) Lambda {
	if len(as) == 0 {
		panic("no argument for function lambda")
	}

	return newLambda(vs, as, e, t)
}

func newLambda(vs []Argument, as []Argument, e Expression, t types.Type) Lambda {
	return Lambda{
		NewLambdaDeclaration(
			argumentsToTypes(vs),
			argumentsToTypes(as),
			t,
		),
		argumentsToNames(vs),
		argumentsToNames(as),
		e,
	}
}

// ArgumentNames returns argument names.
func (l Lambda) ArgumentNames() []string {
	return l.arguments
}

// Body returns a body expression.
func (l Lambda) Body() Expression {
	return l.body
}

// FreeVariableNames returns free varriable names.
func (l Lambda) FreeVariableNames() []string {
	return l.freeVariables
}

// ToDeclaration returns a lambda declaration.
func (l Lambda) ToDeclaration() LambdaDeclaration {
	return l.LambdaDeclaration
}

// ClearFreeVariables clears free variables.
func (l Lambda) ClearFreeVariables() Lambda {
	return Lambda{
		NewLambdaDeclaration(nil, l.ArgumentTypes(), l.ResultType()),
		nil,
		l.arguments,
		l.body,
	}
}

// VisitExpressions visits expressions.
func (l Lambda) VisitExpressions(f func(Expression) error) error {
	return l.body.VisitExpressions(f)
}

// ConvertTypes converts types.
func (l Lambda) ConvertTypes(f func(types.Type) types.Type) Lambda {
	return Lambda{
		l.LambdaDeclaration.ConvertTypes(f),
		l.freeVariables,
		l.arguments,
		l.body.ConvertTypes(f),
	}
}

// RenameVariables renames variables.
func (l Lambda) RenameVariables(vs map[string]string) Lambda {
	fvs := make([]string, 0, len(l.freeVariables))

	for _, s := range l.freeVariables {
		fvs = append(fvs, replaceVariable(s, vs))
	}

	ss := make([]string, 0, len(l.arguments))

	for _, s := range l.arguments {
		ss = append(ss, s)
	}

	return Lambda{
		l.LambdaDeclaration,
		fvs,
		l.arguments,
		l.body.RenameVariables(removeVariables(vs, ss...)),
	}
}

func argumentsToNames(as []Argument) []string {
	ss := make([]string, 0, len(as))

	for _, a := range as {
		ss = append(ss, a.Name())
	}

	// Normalize empty slice representations for tests.
	if len(ss) == 0 {
		ss = nil
	}

	return ss
}

func argumentsToTypes(as []Argument) []types.Type {
	ts := make([]types.Type, 0, len(as))

	for _, a := range as {
		ts = append(ts, a.Type())
	}

	// Normalize empty slice representations for tests.
	if len(ts) == 0 {
		ts = nil
	}

	return ts
}
