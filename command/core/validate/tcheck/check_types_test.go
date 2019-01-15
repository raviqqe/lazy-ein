package tcheck_test

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/core/validate/tcheck"
	"github.com/stretchr/testify/assert"
)

var algebraicType = types.NewAlgebraic(
	types.NewConstructor(types.NewFloat64()),
	types.NewConstructor(),
)

var algebraicLambda = ast.NewVariableLambda(
	nil,
	true,
	ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 1), nil),
	algebraicType,
)

func TestCheckTypes(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		// Empty modules
		nil,
		// Global variables
		{
			ast.NewBind("x", algebraicLambda),
		},
		// Function applications
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewFloat64(),
				),
			),
			ast.NewBind(
				"g",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("f"), []ast.Atom{ast.NewFloat64(42)}),
					types.NewFloat64(),
				),
			),
		},
		// Constructor applications
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(
						ast.NewConstructor(algebraicType, 0),
						[]ast.Atom{ast.NewFloat64(42)},
					),
					algebraicType,
				),
			),
		},
		// Let expressions
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewLet(
						[]ast.Bind{ast.NewBind("a", algebraicLambda)},
						ast.NewFunctionApplication(ast.NewVariable("a"), nil),
					),
					types.NewBoxed(algebraicType),
				),
			),
		},
		// Primitive operations
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewPrimitiveOperation(
						ast.AddFloat64,
						[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
					),
					types.NewFloat64(),
				),
			),
		},
		// Primitive case expressions
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewPrimitiveCase(
						ast.NewFloat64(42),
						types.NewFloat64(),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(0), ast.NewFloat64(1)),
							ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(2049)),
						},
						ast.NewDefaultAlternative(
							"y",
							ast.NewFunctionApplication(ast.NewVariable("y"), nil),
						),
					),
					types.NewFloat64(),
				),
			),
		},
		// Primitive case expressions composed only of default alternatives
		{
			ast.NewBind(
				"x",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewPrimitiveCase(
						ast.NewFloat64(42),
						types.NewFloat64(),
						nil,
						ast.NewDefaultAlternative(
							"y",
							ast.NewFunctionApplication(ast.NewVariable("y"), nil),
						),
					),
					types.NewFloat64(),
				),
			),
		},
		// Primitive case expressions without default alternatives
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewPrimitiveCaseWithoutDefault(
						ast.NewFloat64(42),
						types.NewFloat64(),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(42)),
						},
					),
					types.NewFloat64(),
				),
			),
		},
		// Algebraic case expressions
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewAlgebraicCaseWithoutDefault(
						ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 1), nil),
						[]ast.AlgebraicAlternative{
							ast.NewAlgebraicAlternative(
								ast.NewConstructor(algebraicType, 1),
								nil,
								ast.NewFloat64(42),
							),
							ast.NewAlgebraicAlternative(
								ast.NewConstructor(algebraicType, 0),
								[]string{"y"},
								ast.NewFunctionApplication(ast.NewVariable("y"), nil),
							),
						},
					),
					types.NewFloat64(),
				),
			),
		},
		// Algebraic case expressions with boxed arguments
		{
			ast.NewBind("a", algebraicLambda),
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewAlgebraicCaseWithoutDefault(
						ast.NewFunctionApplication(ast.NewVariable("a"), nil),
						[]ast.AlgebraicAlternative{
							ast.NewAlgebraicAlternative(
								ast.NewConstructor(algebraicType, 1),
								nil,
								ast.NewFloat64(42),
							),
						},
					),
					types.NewFloat64(),
				),
			),
		},
		// Algebraic case expressions with default alternatives
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewAlgebraicCase(
						ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 1), nil),
						[]ast.AlgebraicAlternative{
							ast.NewAlgebraicAlternative(
								ast.NewConstructor(algebraicType, 1),
								nil,
								ast.NewFloat64(42),
							),
						},
						ast.NewDefaultAlternative("z", ast.NewFloat64(42)),
					),
					types.NewFloat64(),
				),
			),
		},
		// Algebraic case expressions composed only of default alternatives
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewAlgebraicCase(
						ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 1), nil),
						nil,
						ast.NewDefaultAlternative("y", ast.NewFloat64(42)),
					),
					types.NewFloat64(),
				),
			),
		},
	} {
		assert.Nil(t, tcheck.CheckTypes(ast.NewModule("", bs)))
	}
}

func TestCheckTypesError(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		// Atoms
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFloat64(42),
					algebraicType,
				),
			),
		},
		// Unknown variables
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
					types.NewBoxed(algebraicType),
				),
			),
		},
		// Non-function calls
		{
			ast.NewBind("f", algebraicLambda),
			ast.NewBind(
				"g",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("f"), []ast.Atom{ast.NewFloat64(42)}),
					algebraicType,
				),
			),
		},
		// Wrong function argument types
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewFloat64(),
				),
			),
			ast.NewBind("a", algebraicLambda),
			ast.NewBind(
				"g",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("f"), []ast.Atom{ast.NewVariable("a")}),
					types.NewFloat64(),
				),
			),
		},
		// Wrong number of function arguments
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewFloat64(),
				),
			),
			ast.NewBind(
				"x",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(
						ast.NewVariable("f"),
						[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
					),
					types.NewFloat64(),
				),
			),
		},
		// Wrong constructor argument types
		{
			ast.NewBind("a", algebraicLambda),
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(
						ast.NewConstructor(algebraicType, 0),
						[]ast.Atom{ast.NewVariable("a")},
					),
					algebraicType,
				),
			),
		},
		// Wrong number of constructor arguments
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 0), nil),
					algebraicType,
				),
			),
		},
		// Wrong primitive operation argument types
		{
			ast.NewBind("a", algebraicLambda),
			ast.NewBind(
				"_",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("_", types.NewFloat64())},
					ast.NewPrimitiveOperation(
						ast.AddFloat64,
						[]ast.Atom{ast.NewFloat64(42), ast.NewVariable("a")},
					),
					types.NewFloat64(),
				),
			),
		},
		// Inconsistent alternative expression types in primitive cases
		{
			ast.NewBind("a", algebraicLambda),
			ast.NewBind(
				"_",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("_", types.NewFloat64())},
					ast.NewPrimitiveCaseWithoutDefault(
						ast.NewFloat64(42),
						types.NewFloat64(),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(0), ast.NewFloat64(1)),
							ast.NewPrimitiveAlternative(
								ast.NewFloat64(42),
								ast.NewFunctionApplication(ast.NewVariable("a"), nil),
							),
						},
					),
					types.NewFloat64(),
				),
			),
		},
	} {
		assert.Error(t, tcheck.CheckTypes(ast.NewModule("", bs)))
	}
}
