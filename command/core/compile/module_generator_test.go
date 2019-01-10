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
	newModuleGenerator(llvm.NewModule("foo"), ast.NewModule("foo", nil))
}

func TestNewModuleGeneratorWithAlgebraicTypes(t *testing.T) {
	for _, b := range []ast.Bind{
		ast.NewBind(
			"foo",
			ast.NewLambda(
				nil,
				true,
				nil,
				ast.NewConstructorApplication("constructor0", []ast.Atom{ast.NewFloat64(42)}),
				types.NewAlgebraic(
					[]types.Constructor{
						types.NewConstructor("constructor0", []types.Type{types.NewFloat64()}),
					},
				),
			),
		),
		ast.NewBind(
			"foo",
			ast.NewLambda(
				nil,
				true,
				nil,
				ast.NewConstructorApplication("constructor0", []ast.Atom{ast.NewFloat64(42)}),
				types.NewAlgebraic(
					[]types.Constructor{
						types.NewConstructor("constructor0", []types.Type{types.NewFloat64()}),
						types.NewConstructor("constructor1", []types.Type{types.NewFloat64(), types.NewFloat64()}),
					},
				),
			),
		),
	} {
		newModuleGenerator(llvm.NewModule("foo"), ast.NewModule("foo", []ast.Bind{b}))
	}
}

func TestModuleGeneratorGenerate(t *testing.T) {
	for _, bs := range [][]ast.Bind{
		// No bind statements
		nil,
		// Global variables
		{
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
		},
		// Functions returning unboxed values
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewFloat64(),
				),
			),
		},
		// Functions returning boxed values
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Function applications
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
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
				"foo",
				ast.NewLambda(
					nil,
					false,
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
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
			ast.NewBind(
				"bar",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewFunctionApplication(ast.NewVariable("foo"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Functions referencing global variables
		{
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
			ast.NewBind(
				"bar",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("foo"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Primitive operations
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
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
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								ast.NewLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									true,
									nil,
									ast.NewFunctionApplication(ast.NewVariable("x"), nil),
									types.NewFloat64(),
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("y"), nil),
					),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Let expressions of functions with free variables
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"bar",
								ast.NewLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									false,
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
							ast.NewVariable("bar"),
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
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"y",
								ast.NewLambda(
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									true,
									nil,
									ast.NewFunctionApplication(ast.NewVariable("x"), nil),
									types.NewFloat64(),
								),
							),
							ast.NewBind(
								"z",
								ast.NewLambda(
									[]ast.Argument{ast.NewArgument("y", types.NewBoxed(types.NewFloat64()))},
									true,
									nil,
									ast.NewFunctionApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(types.NewFloat64()),
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("z"), nil),
					),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Recursive global variables
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewFunctionApplication(ast.NewVariable("foo"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Mutually recursive global variables
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewFunctionApplication(ast.NewVariable("bar"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
			ast.NewBind(
				"bar",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewFunctionApplication(ast.NewVariable("foo"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
		// Recursive functions
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
					ast.NewFunctionApplication(ast.NewVariable("foo"), []ast.Atom{ast.NewVariable("x")}),
					types.NewFloat64(),
				),
			),
		},
		// Recursive let expressions
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"bar",
								ast.NewLambda(
									[]ast.Argument{
										ast.NewArgument(
											"bar",
											types.NewFunction([]types.Type{types.NewFloat64()}, types.NewFloat64())),
									},
									false,
									[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
									ast.NewFunctionApplication(
										ast.NewVariable("bar"),
										[]ast.Atom{ast.NewVariable("x")},
									),
									types.NewFloat64(),
								),
							),
						},
						ast.NewFunctionApplication(
							ast.NewVariable("bar"),
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
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
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
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewPrimitiveCase(
						ast.NewFunctionApplication(ast.NewVariable("x"), nil),
						types.NewBoxed(types.NewFloat64()),
						[]ast.PrimitiveAlternative{
							ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(0)),
						},
						ast.NewDefaultAlternative("x", ast.NewFunctionApplication(ast.NewVariable("x"), nil)),
					),
					types.NewFloat64(),
				),
			),
		},
		// Primitive case expressions without default cases
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
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
		g, err := newModuleGenerator(llvm.NewModule("foo"), ast.NewModule("foo", bs))
		assert.Nil(t, err)
		assert.Nil(t, g.Generate(bs))
	}
}

func TestModuleGeneratorGenerateWithGlobalFunctionsReturningBoxedValues(t *testing.T) {
	m := ast.NewModule(
		"foo",
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewFunctionApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
	)
	mm := llvm.NewModule(m.Name())

	g, err := newModuleGenerator(mm, m)
	assert.Nil(t, err)

	err = g.Generate(m.Binds())
	assert.Nil(t, err)

	assert.Equal(
		t,
		llvm.PointerTypeKind,
		mm.NamedFunction(names.ToEntry("foo")).Type().ElementType().ReturnType().TypeKind(),
	)
}

func TestModuleGeneratorGenerateWithLocalFunctionsReturningBoxedValues(t *testing.T) {
	m := ast.NewModule(
		"foo",
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewLet(
						[]ast.Bind{
							ast.NewBind(
								"bar",
								ast.NewLambda(
									nil,
									false,
									[]ast.Argument{ast.NewArgument("y", types.NewBoxed(types.NewFloat64()))},
									ast.NewFunctionApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(types.NewFloat64()),
								),
							),
						},
						ast.NewFunctionApplication(ast.NewVariable("bar"), []ast.Atom{ast.NewVariable("x")}),
					),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
	)
	mm := llvm.NewModule(m.Name())

	g, err := newModuleGenerator(mm, m)
	assert.Nil(t, err)

	err = g.Generate(m.Binds())
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
	tt := types.NewAlgebraic(
		[]types.Constructor{types.NewConstructor("constructor", []types.Type{types.NewFloat64()})},
	)

	m := ast.NewModule(
		"foo",
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewConstructorApplication("constructor", []ast.Atom{ast.NewFloat64(42)}),
					tt,
				),
			),
		},
	)

	g, err := newModuleGenerator(llvm.NewModule(m.Name()), m)
	assert.Nil(t, err)
	assert.Nil(t, g.Generate(m.Binds()))
}

func TestModuleGeneratorGenerateWithAlgebraicTypesOfMultipleConstructors(t *testing.T) {
	tt := types.NewAlgebraic(
		[]types.Constructor{
			types.NewConstructor("constructor0", []types.Type{types.NewFloat64()}),
			types.NewConstructor("constructor1", []types.Type{types.NewFloat64(), types.NewFloat64()}),
		},
	)

	m := ast.NewModule(
		"foo",
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewConstructorApplication(
						"constructor1",
						[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
					),
					tt,
				),
			),
		},
	)

	g, err := newModuleGenerator(llvm.NewModule(m.Name()), m)
	assert.Nil(t, err)
	assert.Nil(t, g.Generate(m.Binds()))
}

func TestModuleGeneratorGenerateWithAlgebraicCaseExpressions(t *testing.T) {
	tt0 := types.NewAlgebraic(
		[]types.Constructor{types.NewConstructor("constructor", []types.Type{types.NewFloat64()})},
	)
	tt1 := types.NewAlgebraic(
		[]types.Constructor{
			types.NewConstructor("constructor0", []types.Type{types.NewFloat64()}),
			types.NewConstructor("constructor1", []types.Type{types.NewFloat64(), types.NewFloat64()}),
		},
	)

	for _, c := range []ast.Case{
		ast.NewAlgebraicCaseWithoutDefault(
			ast.NewConstructorApplication("constructor", []ast.Atom{ast.NewFloat64(42)}),
			tt0,
			[]ast.AlgebraicAlternative{
				ast.NewAlgebraicAlternative(
					"constructor",
					[]string{"y"},
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
				),
			},
		),
		ast.NewAlgebraicCase(
			ast.NewConstructorApplication("constructor1", []ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)}),
			tt1,
			[]ast.AlgebraicAlternative{
				ast.NewAlgebraicAlternative(
					"constructor1",
					[]string{"x", "y"},
					ast.NewFunctionApplication(ast.NewVariable("y"), nil),
				),
			},
			ast.NewDefaultAlternative("x", ast.NewFloat64(42)),
		),
	} {
		m := ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"foo",
					ast.NewLambda(nil, true, nil, c, types.NewFloat64()),
				),
			},
		)

		g, err := newModuleGenerator(llvm.NewModule(m.Name()), m)
		assert.Nil(t, err)
		assert.Nil(t, g.Generate(m.Binds()))
	}
}
