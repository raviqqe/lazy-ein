package ast

import "github.com/ein-lang/ein/command/core/types"

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

// ConvertTypes converts types.
func (b Bind) ConvertTypes(f func(types.Type) types.Type) Bind {
	return Bind{b.name, b.lambda.ConvertTypes(f)}
}
