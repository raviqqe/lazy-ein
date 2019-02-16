package compile

import (
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator()
}

func TestModuleGeneratorInitializeWithAlgebraicTypes(t *testing.T) {
	tt0 := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))
	tt1 := types.NewAlgebraic(
		types.NewConstructor(types.NewFloat64()),
		types.NewConstructor(types.NewFloat64(), types.NewFloat64()),
	)

	for _, b := range []ast.Bind{
		ast.NewBind(
			"foo",
			ast.NewVariableLambda(
				nil,
				true,
				ast.NewConstructorApplication(
					ast.NewConstructor(tt0, 0),
					[]ast.Atom{ast.NewFloat64(42)},
				),
				tt0,
			),
		),
		ast.NewBind(
			"foo",
			ast.NewVariableLambda(
				nil,
				true,
				ast.NewConstructorApplication(
					ast.NewConstructor(tt1, 0),
					[]ast.Atom{ast.NewFloat64(42)},
				),
				tt1,
			),
		),
	} {
		assert.Nil(t, newModuleGenerator().initialize(ast.NewModule("foo", nil, []ast.Bind{b})))
	}
}

func TestModuleGeneratorGenerate(t *testing.T) {
	a0 := types.NewAlgebraic(types.NewConstructor())
	a1 := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))
	l0 := ast.NewVariableLambda(
		nil,
		true,
		ast.NewConstructorApplication(ast.NewConstructor(a0, 0), nil),
		a0,
	)

	for _, bs := range [][]ast.Bind{
		// No bind statements
		nil,
		// Global variables
		{
			ast.NewBind("x", l0),
		},
		// Functions returning unboxed values
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
		},
		// Functions returning boxed values
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(a0))},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a0),
				),
			),
		},
		// Function applications
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{
						ast.NewArgument(
							"x",
							types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
						),
					},
					ast.NewFunctionApplication(ast.NewVariable("x"), []ast.Atom{ast.NewFloat64(42)}),
					types.NewFloat64(),
				),
			),
		},
		// Function applications with passed arguments
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{
						ast.NewArgument(
							"f",
							types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64()),
						),
						ast.NewArgument(
							"x",
							types.NewFloat64(),
						),
					},
					ast.NewFunctionApplication(ast.NewVariable("f"), []ast.Atom{ast.NewVariable("x")}),
					types.NewFloat64(),
				),
			),
		},
		// Global variables referencing others
		{
			ast.NewBind("x", l0),
			ast.NewBind(
				"y",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a0),
				),
			),
		},
		// Functions referencing global variables
		{
			ast.NewBind("x", l0),
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("y", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a0),
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
		// Let expressions
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								ast.NewVariableLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									true,
									ast.NewConstructorApplication(
										ast.NewConstructor(a1, 0),
										[]ast.Atom{ast.NewVariable("x")},
									),
									a1,
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("y"), nil),
					),
					types.NewBoxed(a1),
				),
			),
		},
		// Let expressions of functions with free variables
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"g",
								ast.NewFunctionLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									[]ast.Argument{ast.NewArgument("y", types.NewFloat64())},
									ast.NewPrimitiveOperation(
										ast.AddFloat64,
										[]ast.Atom{ast.NewVariable("x"), ast.NewVariable("y")},
									),
									types.NewFloat64(),
								),
							),
						},
						ast.NewFunctionApplication(
							ast.NewVariable("g"),
							[]ast.Atom{ast.NewFloat64(42)},
						),
					),
					types.NewFloat64(),
				),
			),
		},
		// Let expressions referencing other bound names
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								ast.NewVariableLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									true,
									ast.NewConstructorApplication(
										ast.NewConstructor(a1, 0),
										[]ast.Atom{ast.NewVariable("x")},
									),
									a1,
								),
							),
							ast.NewBind(
								"z",
								ast.NewVariableLambda(
									[]ast.Argument{ast.NewArgument("y", types.NewBoxed(a1))},
									true,
									ast.NewFunctionApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(a1),
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("z"), nil),
					),
					types.NewBoxed(a1),
				),
			),
		},
		// Recursive global variables
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a0),
				),
			),
		},
		// Mutually recursive global variables
		{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
					types.NewBoxed(a0),
				),
			),
			ast.NewBind(
				"y",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a0),
				),
			),
		},
		// Recursive functions
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("f"), []ast.Atom{ast.NewVariable("x")}),
					types.NewFloat64(),
				),
			),
		},
		// Recursive let expressions
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("_", types.NewBoxed(a0))},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"g",
								ast.NewFunctionLambda(
									[]ast.Argument{
										ast.NewArgument(
											"g",
											types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64())),
									},
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									ast.NewFunctionApplication(
										ast.NewVariable("g"),
										[]ast.Atom{ast.NewVariable("x")},
									),
									types.NewFloat64(),
								),
							),
						},
						ast.NewFunctionApplication(
							ast.NewVariable("g"),
							[]ast.Atom{ast.NewFloat64(42)},
						),
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
					[]ast.Argument{ast.NewArgument("_", types.NewBoxed(a0))},
					ast.NewPrimitiveCase(
						ast.NewFloat64(42),
						types.NewFloat64(),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(0)),
							ast.NewPrimitiveAlternative(ast.NewFloat64(2049), ast.NewFloat64(0)),
						},
						ast.NewDefaultAlternative("x", ast.NewFunctionApplication(ast.NewVariable("x"), nil)),
					),
					types.NewFloat64(),
				),
			),
		},
		// Case expressions unboxing arguments
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(a0))},
					ast.NewAlgebraicCase(
						ast.NewFunctionApplication(ast.NewVariable("x"), nil),
						[]ast.AlgebraicAlternative{
							ast.NewAlgebraicAlternative(ast.NewConstructor(a0, 0), nil, ast.NewFloat64(42)),
						},
						ast.NewDefaultAlternative("x", ast.NewFloat64(42)),
					),
					types.NewFloat64(),
				),
			),
		},
		// Primitive case expressions without default cases
		{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("_", types.NewFloat64())},
					ast.NewPrimitiveCaseWithoutDefault(
						ast.NewFloat64(42),
						types.NewFloat64(),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(0)),
						},
					),
					types.NewFloat64(),
				),
			),
		},
	} {
		_, err := newModuleGenerator().Generate(ast.NewModule("foo", nil, bs))
		assert.Nil(t, err)
	}
}

func TestModuleGeneratorGenerateWithGlobalFunctionsReturningBoxedValues(t *testing.T) {
	a := types.NewAlgebraic(types.NewConstructor())

	m := ast.NewModule(
		"foo",
		nil,
		[]ast.Bind{
			ast.NewBind(
				"f",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(a))},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(a),
				),
			),
		},
	)

	mm, err := newModuleGenerator().Generate(m)
	assert.Nil(t, err)

	assert.Equal(
		t,
		llvm.PointerTypeKind,
		mm.NamedFunction(names.ToEntry("f")).Type().ElementType().ReturnType().TypeKind(),
	)
}

func TestModuleGeneratorGenerateWithLocalFunctionsReturningBoxedValues(t *testing.T) {
	a := types.NewAlgebraic(types.NewConstructor())

	m := ast.NewModule(
		"foo",
		nil,
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewFunctionLambda(
					nil,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(a))},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"bar",
								ast.NewFunctionLambda(
									nil,
									[]ast.Argument{ast.NewArgument("y", types.NewBoxed(a))},
									ast.NewFunctionApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(a),
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("bar"), []ast.Atom{ast.NewVariable("x")}),
					),
					types.NewBoxed(a),
				),
			),
		},
	)

	mm, err := newModuleGenerator().Generate(m)
	assert.Nil(t, err)

	assert.Equal(
		t,
		llvm.PointerTypeKind,
		mm.NamedFunction(
			names.ToEntry(names.NewNameGenerator(names.ToEntry("foo")).Generate("bar")),
		).Type().ElementType().ReturnType().TypeKind(),
	)
}

func TestModuleGeneratorGenerateWithAlgebraicTypes(t *testing.T) {
	tt := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))

	m := ast.NewModule(
		"foo",
		nil,
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(
						ast.NewConstructor(tt, 0),
						[]ast.Atom{ast.NewFloat64(42)},
					),
					tt,
				),
			),
		},
	)

	_, err := newModuleGenerator().Generate(m)
	assert.Nil(t, err)
}

func TestModuleGeneratorGenerateWithAlgebraicTypesOfMultipleConstructors(t *testing.T) {
	tt := types.NewAlgebraic(
		types.NewConstructor(types.NewFloat64()),
		types.NewConstructor(types.NewFloat64(), types.NewFloat64()),
	)

	m := ast.NewModule(
		"foo",
		nil,
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(
						ast.NewConstructor(tt, 1),
						[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
					),
					tt,
				),
			),
		},
	)

	_, err := newModuleGenerator().Generate(m)
	assert.Nil(t, err)
}

func TestModuleGeneratorGenerateWithAlgebraicCaseExpressions(t *testing.T) {
	tt0 := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))
	tt1 := types.NewAlgebraic(
		types.NewConstructor(types.NewFloat64()),
		types.NewConstructor(types.NewFloat64(), types.NewFloat64()),
	)

	for _, c := range []ast.Case{
		ast.NewAlgebraicCaseWithoutDefault(
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							nil,
							true,
							ast.NewConstructorApplication(
								ast.NewConstructor(tt0, 0),
								[]ast.Atom{ast.NewFloat64(42)},
							),
							tt0,
						),
					),
				},
				ast.NewFunctionApplication(ast.NewVariable("x"), nil),
			),
			[]ast.AlgebraicAlternative{
				ast.NewAlgebraicAlternative(
					ast.NewConstructor(tt0, 0),
					[]string{"y"},
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
				),
			},
		),
		ast.NewAlgebraicCase(
			ast.NewLet(
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							nil,
							true,
							ast.NewConstructorApplication(
								ast.NewConstructor(tt1, 1),
								[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
							),
							tt1,
						),
					),
				},
				ast.NewFunctionApplication(ast.NewVariable("x"), nil),
			),
			[]ast.AlgebraicAlternative{
				ast.NewAlgebraicAlternative(
					ast.NewConstructor(tt1, 1),
					[]string{"x", "y"},
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
				),
			},
			ast.NewDefaultAlternative("x", ast.NewFloat64(42)),
		),
	} {
		m := ast.NewModule(
			"f",
			nil,
			[]ast.Bind{
				ast.NewBind(
					"foo",
					ast.NewFunctionLambda(
						nil,
						[]ast.Argument{ast.NewArgument("_", types.NewFloat64())},
						c,
						types.NewFloat64(),
					),
				),
			},
		)

		_, err := newModuleGenerator().Generate(m)
		assert.Nil(t, err)
	}
}

func TestModuleGeneratorGenerateWithDeclarations(t *testing.T) {
	a := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))
	tt := types.NewBoxed(a)

	m := ast.NewModule(
		"foo",
		[]ast.Bind{
			ast.NewBind(
				"x",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewConstructorApplication(
						ast.NewConstructor(a, 0),
						[]ast.Atom{ast.NewFloat64(42)},
					),
					tt,
				),
			),
		},
		[]ast.Bind{
			ast.NewBind(
				"y",
				ast.NewVariableLambda(
					nil,
					true,
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					tt,
				),
			),
		},
	)

	_, err := newModuleGenerator().Generate(m)
	assert.Nil(t, err)
}
