package ast

// Application is a function application
type Application struct {
	function  Variable
	arguments []Atom
}

// NewApplication creates a new application.
func NewApplication(v Variable, as []Atom) Application {
	return Application{v, as}
}

// Function returns a function.
func (a Application) Function() Variable {
	return a.function
}

// Arguments returns arguments.
func (a Application) Arguments() []Atom {
	return a.arguments
}

func (a Application) isExpression() {}
