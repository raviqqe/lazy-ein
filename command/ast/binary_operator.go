package ast

// BinaryOperator is a binary operator.
type BinaryOperator string

const (
	// Add is an addition operator.
	Add BinaryOperator = "+"
	// Subtract is a subtraction operator.
	Subtract = "-"
	// Multiply is a multiplication operator.
	Multiply = "*"
	// Divide is a division operator.
	Divide = "/"
)

// Priority returns operator priority.
func (o BinaryOperator) Priority() int {
	switch o {
	case Add:
		return 1
	case Subtract:
		return 1
	case Multiply:
		return 2
	case Divide:
		return 2
	}

	panic("unreachable")
}
