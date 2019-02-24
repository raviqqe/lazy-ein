package ast

import "github.com/ein-lang/ein/command/core/types"

// PrimitiveOperator is a primitive operator.
type PrimitiveOperator string

const (
	// AddFloat64 is a primitive operator.
	AddFloat64 PrimitiveOperator = "+"
	// SubtractFloat64 is a primitive operator.
	SubtractFloat64 = "-"
	// MultiplyFloat64 is a primitive operator.
	MultiplyFloat64 = "*"
	// DivideFloat64 is a primitive operator.
	DivideFloat64 = "/"
)

// PrimitiveOperation is a saturated primitive operation.
type PrimitiveOperation struct {
	primitiveOperator PrimitiveOperator
	arguments         []Atom
}

// NewPrimitiveOperation creates a primitive operation.
func NewPrimitiveOperation(p PrimitiveOperator, as []Atom) PrimitiveOperation {
	return PrimitiveOperation{p, as}
}

// PrimitiveOperator returns a primitive operator.
func (o PrimitiveOperation) PrimitiveOperator() PrimitiveOperator {
	return o.primitiveOperator
}

// Arguments returns arguments.
func (o PrimitiveOperation) Arguments() []Atom {
	return o.arguments
}

// VisitExpressions visits expressions.
func (o PrimitiveOperation) VisitExpressions(f func(Expression) error) error {
	return f(o)
}

// ConvertTypes converts types.
func (o PrimitiveOperation) ConvertTypes(func(types.Type) types.Type) Expression {
	return o
}

// RenameVariables renames variables.
func (o PrimitiveOperation) RenameVariables(vs map[string]string) Expression {
	as := make([]Atom, 0, len(o.arguments))

	for _, a := range o.arguments {
		as = append(as, a.RenameVariablesInAtom(vs))
	}

	return PrimitiveOperation{o.primitiveOperator, as}
}

func (o PrimitiveOperation) isExpression() {}
