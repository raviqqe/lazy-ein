package codegen

import (
	"testing"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/names"
	"github.com/raviqqe/stg/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

func TestNewModuleGenerator(t *testing.T) {
	newModuleGenerator(llvm.NewModule("foo"), nil)
}

func TestNewModuleGeneratorWithConstructorDefinitions(t *testing.T) {
	for _, d := range []ast.ConstructorDefinition{
		ast.NewConstructorDefinition(
			"foo",
			types.NewAlgebraic(
				[]types.Constructor{
					types.NewConstructor([]types.Type{types.NewFloat64()}),
				},
			),
			0,
		),
		ast.NewConstructorDefinition(
			"foo",
			types.NewAlgebraic(
				[]types.Constructor{
					types.NewConstructor([]types.Type{types.NewFloat64()}),
					types.NewConstructor([]types.Type{types.NewFloat64(), types.NewFloat64()}),
				},
			),
			1,
		),
	} {
		newModuleGenerator(llvm.NewModule("foo"), []ast.ConstructorDefinition{d})
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
					ast.NewApplication(ast.NewVariable("x"), nil),
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
					ast.NewApplication(ast.NewVariable("x"), nil),
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
					ast.NewApplication(ast.NewVariable("x"), []ast.Atom{ast.Float64(42)}),
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
					ast.NewApplication(ast.NewVariable("f"), []ast.Atom{ast.NewVariable("x")}),
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
					ast.NewApplication(ast.NewVariable("foo"), nil),
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
									ast.NewApplication(ast.NewVariable("x"), nil),
									types.NewFloat64(),
								),
							),
						},
						ast.NewApplication(ast.NewVariable("y"), nil),
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
						ast.NewApplication(
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
									ast.NewApplication(ast.NewVariable("x"), nil),
									types.NewFloat64(),
								),
							),
							ast.NewBind(
								"z",
								ast.NewLambda(
									[]ast.Argument{ast.NewArgument("y", types.NewBoxed(types.NewFloat64()))},
									true,
									nil,
									ast.NewApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(types.NewFloat64()),
								),
							),
						},
						ast.NewApplication(ast.NewVariable("z"), nil),
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
					ast.NewApplication("foo", nil),
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
					ast.NewApplication("bar", nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
			ast.NewBind(
				"bar",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewApplication("foo", nil),
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
					ast.NewApplication("foo", []ast.Atom{ast.NewVariable("x")}),
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
									ast.NewApplication(
										ast.NewVariable("bar"),
										[]ast.Atom{ast.NewVariable("x")},
									),
									types.NewFloat64(),
								),
							),
						},
						ast.NewApplication(
							ast.NewVariable("bar"),
							[]ast.Atom{ast.NewFloat64(42)},
						),
					),
					types.NewFloat64(),
				),
			),
		},
		// Case expressions
		{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewCase(
						ast.NewFloat64(42),
						[]ast.Alternative{ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(0))},
						ast.NewDefaultAlternative("x", ast.NewApplication(ast.NewVariable("x"), nil)),
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
					ast.NewCase(
						ast.NewApplication(ast.NewVariable("x"), nil),
						[]ast.Alternative{ast.NewPrimitiveAlternative(ast.NewFloat64(42), ast.NewFloat64(0))},
						ast.NewDefaultAlternative("x", ast.NewApplication(ast.NewVariable("x"), nil)),
					),
					types.NewFloat64(),
				),
			),
		},
	} {
		g, err := newModuleGenerator(llvm.NewModule("foo"), nil)
		assert.Nil(t, err)
		assert.Nil(t, g.Generate(bs))
	}
}

func TestModuleGeneratorGenerateWithGlobalFunctionsReturningBoxedValues(t *testing.T) {
	m := llvm.NewModule("foo")
	g, err := newModuleGenerator(m, nil)
	assert.Nil(t, err)

	err = g.Generate(
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					false,
					[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
					ast.NewApplication(ast.NewVariable("x"), nil),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		llvm.PointerTypeKind,
		m.NamedFunction(names.ToEntry("foo")).Type().ElementType().ReturnType().TypeKind(),
	)
}

func TestModuleGeneratorGenerateWithLocalFunctionsReturningBoxedValues(t *testing.T) {
	m := llvm.NewModule("foo")

	g, err := newModuleGenerator(m, nil)
	assert.Nil(t, err)

	err = g.Generate(
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
									ast.NewApplication(ast.NewVariable("y"), nil),
									types.NewBoxed(types.NewFloat64()),
								),
							),
						},
						ast.NewApplication(ast.NewVariable("bar"), []ast.Atom{ast.Variable("x")}),
					),
					types.NewBoxed(types.NewFloat64()),
				),
			),
		},
	)

	assert.Nil(t, err)

	assert.Equal(
		t,
		llvm.PointerTypeKind,
		m.NamedFunction(
			names.ToEntry(names.NewNameGenerator(names.ToEntry("foo")).Generate("bar")),
		).Type().ElementType().ReturnType().TypeKind(),
	)
}

func TestModuleGeneratorGenerateWithAlgebraicTypes(t *testing.T) {
	tt := types.NewAlgebraic(
		[]types.Constructor{types.NewConstructor([]types.Type{types.NewFloat64()})},
	)

	m := ast.NewModule(
		"foo",
		[]ast.ConstructorDefinition{
			ast.NewConstructorDefinition("constructor", tt, 0),
		},
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewConstructor("constructor", []ast.Atom{ast.NewFloat64(42)}),
					tt,
				),
			),
		},
	)

	g, err := newModuleGenerator(llvm.NewModule(m.Name()), m.ConstructorDefinitions())
	assert.Nil(t, err)
	assert.Nil(t, g.Generate(m.Binds()))
}

func TestModuleGeneratorGenerateWithAlgebraicTypesOfMultipleConstructors(t *testing.T) {
	tt := types.NewAlgebraic(
		[]types.Constructor{
			types.NewConstructor([]types.Type{types.NewFloat64()}),
			types.NewConstructor([]types.Type{types.NewFloat64(), types.NewFloat64()}),
		},
	)

	m := ast.NewModule(
		"foo",
		[]ast.ConstructorDefinition{
			ast.NewConstructorDefinition("constructor1", tt, 0),
			ast.NewConstructorDefinition("constructor2", tt, 1),
		},
		[]ast.Bind{
			ast.NewBind(
				"foo",
				ast.NewLambda(
					nil,
					true,
					nil,
					ast.NewConstructor(
						"constructor2",
						[]ast.Atom{ast.NewFloat64(42), ast.NewFloat64(42)},
					),
					tt,
				),
			),
		},
	)

	g, err := newModuleGenerator(llvm.NewModule(m.Name()), m.ConstructorDefinitions())
	assert.Nil(t, err)
	assert.Nil(t, g.Generate(m.Binds()))
}
