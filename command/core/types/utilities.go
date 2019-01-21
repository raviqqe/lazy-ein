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
	return newEqualityChecker().Check(t, tt)
}
