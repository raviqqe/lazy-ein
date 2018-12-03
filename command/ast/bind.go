package ast

import "github.com/raviqqe/jsonxx/command/types"

// Bind is a bind.
type Bind struct {
	name       string
	arguments  []string
	typ        types.Type
	expression Expression
}

// NewBind creates a new bind.
func NewBind(n string, as []string, t types.Type, e Expression) Bind {
	return Bind{n, as, t, e}
}

// Name returns a name.
func (b Bind) Name() string {
	return b.name
}

// Arguments returns arguments.
func (b Bind) Arguments() []string {
	return b.arguments
}

// Type returns a type.
func (b Bind) Type() types.Type {
	return b.typ
}

// Expression returns a expression.
func (b Bind) Expression() Expression {
	return b.expression
}
