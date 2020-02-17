package ast

import "github.com/raviqqe/lazy-ein/command/types"

// BinaryOperation is a binary operation.
type BinaryOperation struct {
	operator BinaryOperator
	lhs, rhs Expression
}

// NewBinaryOperation creates a binary operation.
func NewBinaryOperation(o BinaryOperator, l, r Expression) BinaryOperation {
	return BinaryOperation{o, l, r}
}

// Operator returns a binary operator.
func (o BinaryOperation) Operator() BinaryOperator {
	return o.operator
}

// LHS returns a left-hand side.
func (o BinaryOperation) LHS() Expression {
	return o.lhs
}

// RHS returns a right-hand side.
func (o BinaryOperation) RHS() Expression {
	return o.rhs
}

// ConvertExpressions converts expressions.
func (o BinaryOperation) ConvertExpressions(f func(Expression) Expression) Expression {
	return f(
		NewBinaryOperation(
			o.operator,
			o.lhs.ConvertExpressions(f).(Expression),
			o.rhs.ConvertExpressions(f).(Expression),
		),
	)
}

// VisitTypes visits types.
func (o BinaryOperation) VisitTypes(f func(types.Type) error) error {
	if err := o.lhs.VisitTypes(f); err != nil {
		return err
	}

	return o.rhs.VisitTypes(f)
}

func (BinaryOperation) isExpression() {}
