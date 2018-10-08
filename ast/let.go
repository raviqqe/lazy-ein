package stg

// Let is a let binding.
type Let struct {
	variable   Variable
	expression Expression
}

func (Let) isExpression() {}
