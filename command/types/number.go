package types

import (
	coreast "github.com/raviqqe/lazy-ein/command/core/ast"
	coretypes "github.com/raviqqe/lazy-ein/command/core/types"
	"github.com/raviqqe/lazy-ein/command/debug"
)

// Number is a number type.
type Number struct {
	debugInformation *debug.Information
}

// NewNumber creates a number type.
func NewNumber(i *debug.Information) Number {
	return Number{i}
}

// Unify unifies itself with another type.
func (n Number) Unify(t Type) ([]Equation, error) {
	if _, ok := t.(Number); ok {
		return nil, nil
	}

	return fallbackToVariable(n, t, NewTypeError("not a number", t.DebugInformation()))
}

// SubstituteVariable substitutes type variables.
func (n Number) SubstituteVariable(v Variable, t Type) Type {
	return n
}

// DebugInformation returns debug information.
func (n Number) DebugInformation() *debug.Information {
	return n.debugInformation
}

// ToCore returns a type in the core language.
func (n Number) ToCore() coretypes.Type {
	return coretypes.NewBoxed(
		coretypes.NewAlgebraic(coretypes.NewConstructor(coretypes.NewFloat64())),
	)
}

// CoreConstructor returns a constructor in the core language.
func (n Number) CoreConstructor() coreast.Constructor {
	return coreast.NewConstructor(coretypes.Unbox(n.ToCore()).(coretypes.Algebraic), 0)
}

// VisitTypes visits types.
func (n Number) VisitTypes(f func(Type) error) error {
	return f(n)
}
