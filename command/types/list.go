package types

import (
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// List is a list type.
type List struct {
	element          Type
	debugInformation *debug.Information
}

// NewList creates a list type.
func NewList(e Type, i *debug.Information) List {
	return List{e, i}
}

// Element returns an element type.
func (l List) Element() Type {
	return l.element
}

// Unify unifies itself with another type.
func (l List) Unify(t Type) ([]Equation, error) {
	ll, ok := t.(List)

	if !ok {
		return fallbackToVariable(l, t, NewTypeError("not a list", t.DebugInformation()))
	}

	return l.element.Unify(ll.element)
}

// SubstituteVariable substitutes type variables.
func (l List) SubstituteVariable(v Variable, t Type) Type {
	return NewList(l.element.SubstituteVariable(v, t), l.debugInformation)
}

// DebugInformation returns debug information.
func (l List) DebugInformation() *debug.Information {
	return l.debugInformation
}

// ToCore returns a type in the core language.
func (l List) ToCore() coretypes.Type {
	return coretypes.NewAlgebraic(
		[]coretypes.Constructor{
			coretypes.NewConstructor(
				l.ConsConstructorName(),
				[]coretypes.Type{
					l.element.ToCore(),
					coretypes.NewBoxed(coretypes.NewReference(l.coreName())),
				},
			),
			coretypes.NewConstructor(l.NilConstructorName(), nil),
		},
	)
}

// VisitTypes visits types.
func (l List) VisitTypes(f func(Type) error) error {
	if err := f(l.element); err != nil {
		return err
	}

	return f(l)
}

// ToTypeDefinition returns a type definition.
func (l List) ToTypeDefinition() coreast.TypeDefinition {
	return coreast.NewTypeDefinition(l.coreName(), l.ToCore())
}

func (l List) coreName() string {
	return "$List." + l.element.coreName() + ".$end"
}

// ConsConstructorName returns a constructor name of cons.
func (l List) ConsConstructorName() string {
	return "$Cons." + l.element.coreName() + ".$end"
}

// NilConstructorName returns a constructor name of nil.
func (l List) NilConstructorName() string {
	return "$Nil." + l.element.coreName() + ".$end"
}
