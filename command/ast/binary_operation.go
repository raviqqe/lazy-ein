package ast

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
func (t BinaryOperation) Operator() BinaryOperator {
	return t.operator
}

// LHS returns a left-hand side.
func (t BinaryOperation) LHS() Expression {
	return t.lhs
}

// RHS returns a right-hand side.
func (t BinaryOperation) RHS() Expression {
	return t.rhs
}

// ConvertExpressions visits expressions.
func (t BinaryOperation) ConvertExpressions(f func(Expression) Expression) node {
	return f(
		NewBinaryOperation(
			t.operator,
			t.lhs.ConvertExpressions(f).(Expression),
			t.rhs.ConvertExpressions(f).(Expression),
		),
	)
}

func (BinaryOperation) isExpression() {}
