package types

import "fmt"

// Function is a function type.
type Function struct {
	arguments []Type
	result    Type
}

// NewFunction creates a function type.
func NewFunction(as []Type, r Type) Function {
	return Function{as, r}
}

// Arguments returns arguments.
func (f Function) Arguments() []Type {
	return f.arguments
}

// Result returns arguments.
func (f Function) Result() Type {
	return f.result
}

func (f Function) String() string {
	s := f.arguments[0].String()

	for _, a := range f.arguments[1:] {
		s += "," + a.String()
	}

	return fmt.Sprintf("Function([%v],%v)", s, f.result)
}

func (Function) isType() {}
