package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Unboxed is an unboxed type.
type Unboxed struct {
	content          Type
	debugInformation *debug.Information
}

// NewUnboxed creates a unboxed type.
func NewUnboxed(c Type, i *debug.Information) Unboxed {
	switch c.(type) {
	case Unboxed:
		panic("cannot unbox unboxed types")
	case Function:
		panic("cannot unbox function types")
	}

	return Unboxed{c, i}
}

// Content returns a content type.
func (u Unboxed) Content() Type {
	return u.content
}

// Unify unifies itself with another type.
func (u Unboxed) Unify(t Type) ([]Equation, error) {
	uu, ok := t.(Unboxed)

	if !ok {
		return fallbackToVariable(u, t, NewTypeError("not an unboxed", t.DebugInformation()))
	}

	return u.content.Unify(uu.content)
}

// SubstituteVariable substitutes type variables.
func (u Unboxed) SubstituteVariable(v Variable, t Type) Type {
	return NewUnboxed(u.content.SubstituteVariable(v, t), u.debugInformation)
}

// DebugInformation returns debug information.
func (u Unboxed) DebugInformation() *debug.Information {
	return u.debugInformation
}

// ToCore returns a type in the core language.
func (u Unboxed) ToCore() coretypes.Type {
	return coretypes.Unbox(u.content.ToCore())
}

// VisitTypes visits types.
func (u Unboxed) VisitTypes(f func(Type) error) error {
	if err := f(u.content); err != nil {
		return err
	}

	return f(u)
}

func (u Unboxed) coreName() string {
	return "$Unboxed." + u.content.coreName() + ".$end"
}
