package ast

import (
	"github.com/ein-lang/ein/command/core/types"
)

// LambdaDeclaration is a declaration.
type LambdaDeclaration struct {
	freeVariables []types.Type
	updatable     bool
	arguments     []types.Type
	result        types.Type
}

// NewLambdaDeclaration creates a declaration.
func NewLambdaDeclaration(vs []types.Type, u bool, as []types.Type, r types.Type) LambdaDeclaration {
	return LambdaDeclaration{vs, u, as, r}
}

// FreeVariableTypes returns free variable types.
func (l LambdaDeclaration) FreeVariableTypes() []types.Type {
	return l.freeVariables
}

// ArgumentTypes returns argument types.
func (l LambdaDeclaration) ArgumentTypes() []types.Type {
	return l.arguments
}

// ResultType returns a result type.
func (l LambdaDeclaration) ResultType() types.Type {
	return l.result
}

// IsUpdatable returns true if the lambda form is updatable, or false otherwise.
func (l LambdaDeclaration) IsUpdatable() bool {
	return l.updatable
}

// IsThunk returns true if the lambda form is a thunk, or false otherwise.
func (l LambdaDeclaration) IsThunk() bool {
	return len(l.arguments) == 0
}

// ConvertTypes converts types.
func (l LambdaDeclaration) ConvertTypes(f func(types.Type) types.Type) LambdaDeclaration {
	vs := make([]types.Type, 0, len(l.freeVariables))

	for _, v := range l.freeVariables {
		vs = append(vs, v.ConvertTypes(f))
	}

	as := make([]types.Type, 0, len(l.arguments))

	for _, a := range l.arguments {
		as = append(as, a.ConvertTypes(f))
	}

	return NewLambdaDeclaration(vs, l.updatable, as, l.result.ConvertTypes(f))
}
