package canonicalize

import "github.com/ein-lang/ein/command/core/types"

type typeValidator struct {
	environment []types.Type
}

func newTypeValidator() typeValidator {
	return typeValidator{}
}

func (v typeValidator) Validate(t types.Type) bool {
	switch t := t.(type) {
	case types.Boxed:
		return v.Validate(t.Content())
	case types.Algebraic:
		v = v.pushType(t)

		for _, c := range t.Constructors() {
			for _, e := range c.Elements() {
				if !v.Validate(e) {
					return false
				}
			}
		}

		return true
	case types.Function:
		v = v.pushType(t)

		for _, a := range t.Arguments() {
			if !v.Validate(a) {
				return false
			}
		}

		return v.Validate(t.Result())
	case types.Index:
		return t.Value() < len(v.environment)
	}

	return true
}

func (v typeValidator) pushType(t types.Type) typeValidator {
	return typeValidator{append(v.environment, t)}
}
