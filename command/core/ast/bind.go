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

// VisitExpressions visits expressions.
func (b Bind) VisitExpressions(f func(Expression) error) error {
	return b.lambda.VisitExpressions(f)
}

// ConvertTypes converts types.
func (b Bind) ConvertTypes(f func(types.Type) types.Type) Bind {
	return Bind{b.name, b.lambda.ConvertTypes(f)}
}

// RenameVariables renames variables.
func (b Bind) RenameVariables(vs map[string]string) Bind {
	return Bind{b.name, b.lambda.RenameVariables(vs)}
}
