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
func (l List) ToCore() (coretypes.Type, error) {
	s, err := l.coreName()

	if err != nil {
		return nil, err
	}

	return coretypes.NewNamed(s), nil
}

// ToTypeDefinition returns a type definition.
func (l List) ToTypeDefinition() (coreast.TypeDefinition, error) {
	t, err := l.element.ToCore()

	if err != nil {
		return coreast.TypeDefinition{}, err
	}

	s, err := l.coreName()

	if err != nil {
		return coreast.TypeDefinition{}, err
	}

	return coreast.NewTypeDefinition(
		s,
		coretypes.NewAlgebraic(
			[]coretypes.Constructor{
				coretypes.NewConstructor(
					"$Cons",
					[]coretypes.Type{
						t,
						coretypes.NewBoxed(coretypes.NewNamed(s)),
					},
				),
				coretypes.NewConstructor("$Nil", nil),
			},
		),
	), nil
}

func (l List) coreName() (string, error) {
	s, err := l.element.coreName()

	if err != nil {
		return "", err
	}

	return "$List." + s + ".$end", nil
}
