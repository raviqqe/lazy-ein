package ast

import "github.com/ein-lang/ein/command/core/types"

// Lambda is a lambda form.
type Lambda struct {
	freeVariables []Argument
	updatable     bool
	arguments     []Argument
	body          Expression
	resultType    types.Type
}

// NewVariableLambda creates a lambda form.
func NewVariableLambda(vs []Argument, u bool, e Expression, t types.Bindable) Lambda {
	return Lambda{vs, u, nil, e, t}
}

// NewFunctionLambda creates a lambda form.
func NewFunctionLambda(vs []Argument, as []Argument, e Expression, t types.Type) Lambda {
	if len(as) == 0 {
		panic("no argument for function lambda")
	}

	return Lambda{vs, false, as, e, t}
}

// Type returns a type.
func (l Lambda) Type() types.Type {
	if len(l.arguments) == 0 {
		return types.Box(l.resultType)
	}

	return types.NewFunction(l.ArgumentTypes(), l.resultType)
}

// Arguments returns arguments.
func (l Lambda) Arguments() []Argument {
	return l.arguments
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

// IsThunk returns true if the lambda form is a thunk, or false otherwise.
func (l Lambda) IsThunk() bool {
	return len(l.arguments) == 0
}

// ClearFreeVariables clears free variables.
func (l Lambda) ClearFreeVariables() Lambda {
	return Lambda{nil, l.updatable, l.arguments, l.body, l.resultType}
}

// ConvertTypes converts types.
func (l Lambda) ConvertTypes(f func(types.Type) types.Type) Lambda {
	vs := make([]Argument, 0, len(l.freeVariables))

	for _, v := range l.freeVariables {
		vs = append(vs, v.ConvertTypes(f))
	}

	as := make([]Argument, 0, len(l.arguments))

	for _, a := range l.arguments {
		as = append(as, a.ConvertTypes(f))
	}

	return Lambda{vs, l.updatable, as, l.body.ConvertTypes(f), l.resultType.ConvertTypes(f)}
}

// RenameVariables renames variables.
func (l Lambda) RenameVariables(vs map[string]string) Lambda {
	fvs := make([]Argument, 0, len(l.freeVariables))

	for _, v := range l.freeVariables {
		fvs = append(fvs, v.RenameVariables(vs))
	}

	ss := make([]string, 0, len(l.arguments))

	for _, a := range l.arguments {
		ss = append(ss, a.Name())
	}

	return Lambda{
		fvs,
		l.updatable,
		l.arguments,
		l.body.RenameVariables(removeVariables(vs, ss...)),
		l.resultType,
	}
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
