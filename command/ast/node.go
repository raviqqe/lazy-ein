package ast

type node interface {
	ConvertExpression(func(Expression) Expression) node
}
