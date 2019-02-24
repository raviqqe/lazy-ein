package ast

import "github.com/ein-lang/ein/command/types"

// Application is a function application.
type Application struct {
	function  Expression
	arguments []Expression
}

// NewApplication creates a function application.
func NewApplication(f Expression, as []Expression) Application {
	return Application{f, as}
}

// Function returns a function.
func (a Application) Function() Expression {
	return a.function
}

// Arguments returns arguments.
func (a Application) Arguments() []Expression {
	return a.arguments
}

// ConvertExpressions converts expressions.
func (a Application) ConvertExpressions(f func(Expression) Expression) Expression {
	as := make([]Expression, 0, len(a.arguments))

	for _, a := range a.arguments {
		as = append(as, a.ConvertExpressions(f).(Expression))
	}

	return f(NewApplication(a.function.ConvertExpressions(f).(Expression), as))
}

// VisitTypes visits types.
func (a Application) VisitTypes(f func(types.Type) error) error {
	if err := a.function.VisitTypes(f); err != nil {
		return err
	}

	for _, a := range a.arguments {
		if err := a.VisitTypes(f); err != nil {
			return err
		}
	}

	return nil
}

func (Application) isExpression() {}
