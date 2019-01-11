package ast

import (
	"fmt"

	"github.com/ein-lang/ein/command/core/types"
)

// Constructor is a constructor.
type Constructor struct {
	typ   types.Algebraic
	index int
}

// NewConstructor creates a constructor.
func NewConstructor(t types.Algebraic, i int) Constructor {
	return Constructor{t, i}
}

// AlgebraicType returns an algebraic type.
func (c Constructor) AlgebraicType() types.Algebraic {
	return c.typ
}

// ConstructorType returns a constructor type.
func (c Constructor) ConstructorType() types.Constructor {
	return c.typ.Constructors()[c.index]
}

// Index returns a constructor index.
func (c Constructor) Index() int {
	return c.index
}

// ID returns an ID.
func (c Constructor) ID() string {
	return fmt.Sprintf("%v.%v", c.typ, c.index)
}
