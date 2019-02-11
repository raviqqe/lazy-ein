package types

import "fmt"

// Constructor is a constructor.
type Constructor struct {
	elements []Type
}

// NewConstructor creates a constructor.
func NewConstructor(ts ...Type) Constructor {
	return Constructor{ts}
}

// Elements returns element types.
func (c Constructor) Elements() []Type {
	return c.elements
}

// ConvertTypes converts types.
func (c Constructor) ConvertTypes(f func(Type) Type) Constructor {
	es := []Type(nil) // HACK: Do not use make() for equality check.

	for _, e := range c.elements {
		es = append(es, e.ConvertTypes(f))
	}

	return Constructor{es}
}

func (c Constructor) String() string {
	if len(c.elements) == 0 {
		return "c"
	}

	s := c.elements[0].String()

	for _, e := range c.elements[1:] {
		s += "," + e.String()
	}

	return fmt.Sprintf("c(%v)", s)
}

func (c Constructor) equal(cc Constructor) bool {
	if len(c.elements) != len(cc.elements) {
		return false
	}

	for i, e := range c.elements {
		if !e.equal(cc.elements[i]) {
			return false
		}
	}

	return true
}
