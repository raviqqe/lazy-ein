package types

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
