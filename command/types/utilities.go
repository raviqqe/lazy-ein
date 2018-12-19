package types

// Box boxes a type.
func Box(t Type) Type {
	if u, ok := t.(Unboxed); ok {
		return u.Content()
	}

	return t
}

func fallbackToVariable(t, tt Type, err error) ([]Equation, error) {
	if _, ok := tt.(Variable); ok {
		return []Equation{NewEquation(t, tt)}, nil
	}

	return nil, err
}
