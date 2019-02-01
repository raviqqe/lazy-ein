package ast

import "github.com/ein-lang/ein/command/types"

// List is a list.
type List struct {
	typ      types.Type
	elements []Expression
}

// NewList creates a list.
func NewList(t types.Type, es []Expression) List {
	return List{t, es}
}

// Type returns a type.
func (l List) Type() types.Type {
	return l.typ
}

// Elements returns elements.
func (l List) Elements() []Expression {
	return l.elements
}

// ConvertExpressions visits expressions.
func (l List) ConvertExpressions(f func(Expression) Expression) Node {
	es := make([]Expression, 0, len(l.elements))

	for _, e := range l.elements {
		es = append(es, e.ConvertExpressions(f).(Expression))
	}

	return f(NewList(l.typ, es))
}

// VisitTypes visits types.
func (l List) VisitTypes(f func(types.Type) error) error {
	if err := l.typ.VisitTypes(f); err != nil {
		return err
	}

	for _, e := range l.elements {
		if err := e.VisitTypes(f); err != nil {
			return err
		}
	}

	return nil
}

func (List) isExpression() {}
