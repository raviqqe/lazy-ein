package types

// Box boxes a type.
func Box(t Type) Type {
	if a, ok := t.(Algebraic); ok {
		return NewBoxed(a)
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
	return t.equal(tt)
}
