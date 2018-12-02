package ast

// Bind is a bind.
type Bind struct {
	name       string
	expression Expression
}

// NewBind creates a new bind.
func NewBind(n string, e Expression) Bind {
	return Bind{n, e}
}

// Name returns a name.
func (b Bind) Name() string {
	return b.name
}

// Expression returns a expression.
func (b Bind) Expression() Expression {
	return b.expression
}
