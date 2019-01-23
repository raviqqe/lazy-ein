package types

type validator struct {
	environment []Type
}

func newValidator() validator {
	return validator{}
}

func (v validator) Validate(t Type) bool {
	switch t := t.(type) {
	case Boxed:
		return v.Validate(t.Content())
	case Algebraic:
		v = v.pushType(t)

		for _, c := range t.Constructors() {
			for _, e := range c.Elements() {
				if !v.Validate(e) {
					return false
				}
			}
		}

		return true
	case Function:
		v = v.pushType(t)

		for _, a := range t.Arguments() {
			if !v.Validate(a) {
				return false
			}
		}

		return v.Validate(t.Result())
	case Index:
		return t.Value() < len(v.environment)
	}

	return true
}

func (v validator) pushType(t Type) validator {
	return validator{append(v.environment, t)}
}
