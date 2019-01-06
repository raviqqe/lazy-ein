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

// NewFunction creates a function type.
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

// ArgumentsCount returns a number of arguments.
func (f Function) ArgumentsCount() int {
	if f, ok := f.result.(Function); ok {
		return 1 + f.ArgumentsCount()
	}

	return 1
}

// Unify unifies itself with another type.
func (f Function) Unify(t Type) ([]Equation, error) {
	ff, ok := t.(Function)

	if !ok {
		return fallbackToVariable(f, t, NewTypeError("not a function", t.DebugInformation()))
	}

	es, err := f.argument.Unify(ff.argument)

	if err != nil {
		return nil, err
	}

	ees, err := f.result.Unify(ff.result)

	if err != nil {
		return nil, err
	}

	return append(es, ees...), nil
}

// SubstituteVariable substitutes type variables.
func (f Function) SubstituteVariable(v Variable, t Type) Type {
	return NewFunction(
		f.argument.SubstituteVariable(v, t),
		f.result.SubstituteVariable(v, t),
		f.debugInformation,
	)
}

// DebugInformation returns debug information.
func (f Function) DebugInformation() *debug.Information {
	return f.debugInformation
}

// ToCore returns a type in the core language.
func (f Function) ToCore() (coretypes.Type, error) {
	as := []coretypes.Type{}

	for {
		a, err := f.Argument().ToCore()

		if err != nil {
			return nil, err
		}

		as = append(as, a)
		ff, ok := f.Result().(Function)

		if !ok {
			r, err := f.Result().ToCore()

			if err != nil {
				return nil, err
			}

			return coretypes.NewFunction(as, r), nil
		}

		f = ff
	}
}

// VisitTypes visits types.
func (f Function) VisitTypes(ff func(Type) error) error {
	if err := ff(f.argument); err != nil {
		return err
	} else if err := ff(f.result); err != nil {
		return err
	}

	return ff(f)
}

func (f Function) coreName() (string, error) {
	a, err := f.argument.coreName()

	if err != nil {
		return "", err
	}

	r, err := f.result.coreName()

	if err != nil {
		return "", err
	}

	return "$Function." + a + "." + r + ".$end", nil
}
