package ast

type node interface {
	ConvertExpressions(func(Expression) Expression) node
}
