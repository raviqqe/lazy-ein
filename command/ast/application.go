package ast

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

// ConvertExpression visits expressions.
func (a Application) ConvertExpression(f func(Expression) Expression) node {
	as := make([]Expression, 0, len(a.arguments))

	for _, a := range a.arguments {
		as = append(as, f(a))
	}

	return NewApplication(f(a.function), as)
}

func (Application) isExpression() {}
