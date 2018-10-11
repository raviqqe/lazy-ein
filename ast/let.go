package ast

// Let is a let binding.
type Let struct {
	binds      []Bind
	expression Expression
}

func (Let) isExpression() {}
