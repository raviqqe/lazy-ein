package types

// Equation is a type equation.
type Equation struct {
	left  Type
	right Type
}

// NewEquation creates a type equation.
func NewEquation(t, tt Type) Equation {
	return Equation{t, tt}
}

// Left returns a left side of an equation.
func (e Equation) Left() Type {
	return e.left
}

// Right returns a right side of an equation.
func (e Equation) Right() Type {
	return e.right
}
