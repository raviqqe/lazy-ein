package types

// Constructor is a constructor.
type Constructor struct {
	elements []Type
}

// NewConstructor creates a constructor.
func NewConstructor(es []Type) Constructor {
	return Constructor{es}
}

// Elements returns element types.
func (c Constructor) Elements() []Type {
	return c.elements
}
