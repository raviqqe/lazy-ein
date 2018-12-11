package types

import "github.com/ein-lang/ein/command/debug"

// Unboxed is an unboxed type.
type Unboxed struct {
	content          Type
	debugInformation *debug.Information
}

// NewUnboxed creates a unboxed type.
func NewUnboxed(c Type, i *debug.Information) Unboxed {
	if _, ok := c.(Unboxed); ok {
		panic("cannot unbox unboxed types")
	}

	return Unboxed{c, i}
}

// Content returns a content type.
func (u Unboxed) Content() Type {
	return u.content
}

// Unify unifies itself with another type.
func (u Unboxed) Unify(t Type) error {
	uu, ok := t.(Unboxed)

	if !ok {
		return newTypeError("not an unboxed", t.DebugInformation())
	}

	return u.content.Unify(uu.content)
}

// DebugInformation returns debug information.
func (u Unboxed) DebugInformation() *debug.Information {
	return u.debugInformation
}
