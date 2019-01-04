package types

// Constructor is a constructor.
type Constructor struct {
	name     string
	elements []Type
}

// NewConstructor creates a constructor.
func NewConstructor(n string, es []Type) Constructor {
	return Constructor{n, es}
}

// Name returns a name.
func (c Constructor) Name() string {
	return c.name
}

// Elements returns element types.
func (c Constructor) Elements() []Type {
	return c.elements
}
