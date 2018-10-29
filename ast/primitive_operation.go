package ast

// Primitive is a primitive operation.
type Primitive string

const (
	// AddFloat64 is a primitive operation.
	AddFloat64 Primitive = "+"
	// SubtractFloat64 is a primitive operation.
	SubtractFloat64 = "-"
	// MultiplyFloat64 is a primitive operation.
	MultiplyFloat64 = "*"
	// DivideFloat64 is a primitive operation.
	DivideFloat64 = "/"
)

// PrimitiveOperation is a saturated primitive operation.
type PrimitiveOperation struct {
	primitive Primitive
	arguments []Atom
}

// NewPrimitiveOperation creates a new primitive operation.
func NewPrimitiveOperation(p Primitive, as []Atom) PrimitiveOperation {
	return PrimitiveOperation{p, as}
}

// Primitive returns a primitive.
func (a PrimitiveOperation) Primitive() Primitive {
	return a.primitive
}

// Arguments returns arguments.
func (a PrimitiveOperation) Arguments() []Atom {
	return a.arguments
}

func (a PrimitiveOperation) isExpression() {}
