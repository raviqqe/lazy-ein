package ast

import "github.com/raviqqe/lazy-ein/command/types"

// List is a list.
type List struct {
	typ       types.Type
	arguments []ListArgument
}

// NewList creates a list.
func NewList(t types.Type, as []ListArgument) List {
	return List{t, as}
}

// Type returns a type.
func (l List) Type() types.Type {
	return l.typ
}

// Arguments returns arguments.
func (l List) Arguments() []ListArgument {
	return l.arguments
}

// ConvertExpressions converts expressions.
func (l List) ConvertExpressions(f func(Expression) Expression) Expression {
	as := make([]ListArgument, 0, len(l.arguments))

	for _, a := range l.arguments {
		as = append(as, a.ConvertExpressions(f))
	}

	return f(NewList(l.typ, as))
}

// VisitTypes visits types.
func (l List) VisitTypes(f func(types.Type) error) error {
	if err := l.typ.VisitTypes(f); err != nil {
		return err
	}

	for _, a := range l.arguments {
		if err := a.VisitTypes(f); err != nil {
			return err
		}
	}

	return nil
}

func (List) isExpression() {}
