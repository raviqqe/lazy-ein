package ast

// Bind is a bind statement.
type Bind struct {
	name   string
	lambda Lambda
}

// NewBind creates a bind statement.
func NewBind(n string, l Lambda) Bind {
	return Bind{n, l}
}

// Name returns a name.
func (b Bind) Name() string {
	return b.name
}

// Lambda returns a lambda form.
func (b Bind) Lambda() Lambda {
	return b.lambda
}
