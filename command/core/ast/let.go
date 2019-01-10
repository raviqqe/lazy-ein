package ast

import "github.com/ein-lang/ein/command/core/types"

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

// ConvertTypes converts types.
func (l Let) ConvertTypes(f func(types.Type) types.Type) Expression {
	bs := make([]Bind, 0, len(l.binds))

	for _, b := range l.binds {
		bs = append(bs, b.ConvertTypes(f))
	}

	return Let{bs, l.expression.ConvertTypes(f)}
}

func (Let) isExpression() {}
