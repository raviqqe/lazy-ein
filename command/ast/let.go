package ast

// Let is a let expression.
type Let struct {
	binds      []Bind
	expression Expression
}

// NewLet creates a let expression.
func NewLet(bs []Bind, e Expression) Let {
	return Let{bs, e}
}

// Binds returns binds.
func (l Let) Binds() []Bind {
	return l.binds
}

// Expression returns an expression.
func (l Let) Expression() Expression {
	return l.expression
}

func (Let) isExpression() {}

// ConvertExpressions visits expressions.
func (l Let) ConvertExpressions(f func(Expression) Expression) Node {
	bs := make([]Bind, 0, len(l.binds))

	for _, b := range l.binds {
		bs = append(bs, b.ConvertExpressions(f).(Bind))
	}

	return f(NewLet(bs, l.expression.ConvertExpressions(f).(Expression)))
}
