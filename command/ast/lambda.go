package ast

// Lambda is a lambda.
type Lambda struct {
	arguments  []string
	expression Expression
}

// NewLambda creates a lambda.
func NewLambda(as []string, e Expression) Lambda {
	return Lambda{as, e}
}

// Arguments returns arguments.
func (b Lambda) Arguments() []string {
	return b.arguments
}

// Expression returns an expression.
func (b Lambda) Expression() Expression {
	return b.expression
}

// ConvertExpressions visits expressions.
func (b Lambda) ConvertExpressions(f func(Expression) Expression) node {
	return f(NewLambda(b.arguments, b.expression.ConvertExpressions(f).(Expression)))
}

func (Lambda) isExpression() {}
