package ast

// BinaryOperatorTerm is a binary operator term.
type BinaryOperatorTerm struct {
	operator BinaryOperator
	lhs, rhs Expression
}

// NewBinaryOperatorTerm creates a binary operator term.
func NewBinaryOperatorTerm(o BinaryOperator, l, r Expression) BinaryOperatorTerm {
	return BinaryOperatorTerm{o, l, r}
}

// Operator returns a binary operator.
func (t BinaryOperatorTerm) Operator() BinaryOperator {
	return t.operator
}

// LHS returns a left-hand side.
func (t BinaryOperatorTerm) LHS() Expression {
	return t.lhs
}

// RHS returns a right-hand side.
func (t BinaryOperatorTerm) RHS() Expression {
	return t.rhs
}

// ConvertExpression visits expressions.
func (t BinaryOperatorTerm) ConvertExpression(f func(Expression) Expression) node {
	return f(
		NewBinaryOperatorTerm(
			t.operator,
			t.lhs.ConvertExpression(f).(Expression),
			t.rhs.ConvertExpression(f).(Expression),
		),
	)
}

func (BinaryOperatorTerm) isExpression() {}
