package ast

// Literal is a literal.
type Literal interface {
	Expression
	isLiteral()
}

// IsLiteral checks if an expression is a literal or not.
// This function is added only to fix the megacheck linter error.
func IsLiteral(e Expression) bool {
	_, ok := e.(Literal)
	return ok
}
