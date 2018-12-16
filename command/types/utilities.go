package types

// Box boxes a type.
func Box(t Type) Type {
	if u, ok := t.(Unboxed); ok {
		return u.Content()
	}

	return t
}
