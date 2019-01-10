package types

import "fmt"

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

func (c Constructor) String() string {
	if len(c.elements) == 0 {
		return fmt.Sprintf("Constructor(%v)", c.name)
	}

	s := c.elements[0].String()

	for _, e := range c.elements[1:] {
		s += "," + e.String()
	}

	return fmt.Sprintf("Constructor(%v,[%v])", c.name, s)
}

func (c Constructor) equal(cc Constructor) bool {
	if c.name != cc.name || len(c.elements) != len(cc.elements) {
		return false
	}

	for i, e := range c.elements {
		if !e.equal(cc.elements[i]) {
			return false
		}
	}

	return true
}
