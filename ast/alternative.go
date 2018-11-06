package ast

// Alternative is an alternative.
type Alternative interface {
	Expression() Expression
	isAlternative()
}
