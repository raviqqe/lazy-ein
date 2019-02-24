package ast

import "github.com/ein-lang/ein/command/core/types"

// FunctionApplication is a function application.
type FunctionApplication struct {
	function  Variable
	arguments []Atom
}

// NewFunctionApplication creates an application.
func NewFunctionApplication(v Variable, as []Atom) FunctionApplication {
	return FunctionApplication{v, as}
}

// Function returns a function.
func (a FunctionApplication) Function() Variable {
	return a.function
}

// Arguments returns arguments.
func (a FunctionApplication) Arguments() []Atom {
	return a.arguments
}

// VisitExpressions visits expressions.
func (a FunctionApplication) VisitExpressions(f func(Expression) error) error {
	return f(a)
}

// ConvertTypes converts types.
func (a FunctionApplication) ConvertTypes(func(types.Type) types.Type) Expression {
	return a
}

// RenameVariables renames variables.
func (a FunctionApplication) RenameVariables(vs map[string]string) Expression {
	as := make([]Atom, 0, len(a.arguments))

	for _, a := range a.arguments {
		as = append(as, a.RenameVariablesInAtom(vs))
	}

	return FunctionApplication{a.function.RenameVariablesInAtom(vs).(Variable), as}
}

func (a FunctionApplication) isExpression() {}
