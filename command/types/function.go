package types

import (
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/debug"
)

// Function is a function type.
type Function struct {
	argument         Type
	result           Type
	debugInformation *debug.Information
}

// NewFunction creates a new function type.
func NewFunction(a, r Type, i *debug.Information) Function {
	return Function{a, r, i}
}

// Argument returns an argument type.
func (f Function) Argument() Type {
	return f.argument
}

// Result returns a result type.
func (f Function) Result() Type {
	return f.result
}

// Unify unifies itself with another type.
func (f Function) Unify(t Type) error {
	ff, ok := t.(Function)

	if !ok {
		return newTypeError("not a function", t.DebugInformation())
	} else if err := ff.argument.Unify(f.argument); err != nil {
		return err
	}

	return f.result.Unify(ff.result)
}

// DebugInformation returns debug information.
func (f Function) DebugInformation() *debug.Information {
	return f.debugInformation
}

// ToCore returns a type in the core language.
func (f Function) ToCore() coretypes.Type {
	as := []coretypes.Type{}

	for {
		as = append(as, f.Argument().ToCore())
		t, ok := f.Result().(Function)

		if !ok {
			return coretypes.NewFunction(as, f.Result().ToCore())
		}

		f = t
	}
}
