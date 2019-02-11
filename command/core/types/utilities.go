package types

// Box boxes a type.
func Box(t Type) Type {
	if t, ok := t.(Boxable); ok {
		return NewBoxed(t)
	}

	return t
}

// Unbox converts a type into its unboxed type.
func Unbox(t Type) Type {
	if t, ok := t.(Boxed); ok {
		return t.Content()
	}

	return t
}

// Equal checks type equality.
func Equal(t, tt Type) bool {
	return newEqualityChecker(nil).Check(t, tt)
}

// EqualWithEnvironment checks type equality with environment.
func EqualWithEnvironment(t, tt Type, ts []Type) bool {
	return newEqualityChecker(ts).Check(t, tt)
}

// Validate validates a type.
func Validate(t Type) bool {
	return newValidator().Validate(t)
}

// IsRecursive checks if a type is recursive.
func IsRecursive(t Type) bool {
	return newRecursivityChecker().Check(t)
}

// Unwrap unwraps a recursive type by a level.
func Unwrap(t Type) Type {
	return unwrap(t, nil)
}

func unwrap(t Type, ts []Type) Type {
	switch t := t.(type) {
	case Algebraic:
		ts = append(ts, t)
		cs := make([]Constructor, 0, len(t.Constructors()))

		for _, c := range t.Constructors() {
			es := []Type(nil) // HACK: Do not use make() for equality check.

			for _, e := range c.Elements() {
				es = append(es, unwrap(e, ts))
			}

			cs = append(cs, NewConstructor(es...))
		}

		return NewAlgebraic(cs[0], cs[1:]...)
	case Boxed:
		return NewBoxed(unwrap(t.Content(), ts).(Boxable))
	case Function:
		ts = append(ts, t)
		as := make([]Type, 0, len(t.Arguments()))

		for _, a := range t.Arguments() {
			as = append(as, unwrap(a, ts))
		}

		return NewFunction(as, t.Result())
	case Index:
		if t.Value() == len(ts)-1 {
			return ts[0]
		}
	}

	return t
}
