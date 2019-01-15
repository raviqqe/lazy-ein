package ast

// Case is a case expression.
type Case interface {
	Expression
	Argument() Expression
	DefaultAlternative() (DefaultAlternative, bool)
}
