package compile

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	coreast "github.com/raviqqe/lazy-ein/command/core/ast"
)

func binaryOperatorToPrimitive(o ast.BinaryOperator) coreast.PrimitiveOperator {
	switch o {
	case ast.Add:
		return coreast.AddFloat64
	case ast.Subtract:
		return coreast.SubtractFloat64
	case ast.Multiply:
		return coreast.MultiplyFloat64
	case ast.Divide:
		return coreast.DivideFloat64
	}

	panic("unreachable")
}
