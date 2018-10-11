package ast

// Bind is a bind statement.
type Bind struct {
	name   string
	lambda Lambda
}
