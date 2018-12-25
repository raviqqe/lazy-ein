package compile

import (
	"testing"

	"github.com/ein-lang/ein/command/ast"
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestCompileWithEmptySource(t *testing.T) {
	_, err := Compile(ast.NewModule("", []ast.Bind{}))
	assert.Nil(t, err)
}

func TestCompileWithFunctionApplications(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
				),
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewApplication(
						ast.NewVariable("f"),
						[]ast.Expression{ast.NewNumber(42)},
					),
				),
			},
		),
	)

	assert.Nil(t, err)
}

func TestCompileWithNestedFunctionApplications(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
				),
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewApplication(
						ast.NewVariable("f"),
						[]ast.Expression{
							ast.NewApplication(
								ast.NewVariable("f"),
								[]ast.Expression{ast.NewVariable("x")},
							),
						},
					),
				),
			},
		),
	)

	assert.Nil(t, err)
}

func TestCompileWithDeeplyNestedFunctionApplicationsInLambdaExpressions(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda([]string{"x"}, ast.NewVariable("x")),
				),
				ast.NewBind(
					"g",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda(
						[]string{"x"},
						ast.NewApplication(
							ast.NewVariable("f"),
							[]ast.Expression{
								ast.NewApplication(
									ast.NewVariable("f"),
									[]ast.Expression{
										ast.NewApplication(
											ast.NewVariable("f"),
											[]ast.Expression{ast.NewVariable("x")},
										),
									},
								),
							},
						),
					),
				),
			},
		),
	)

	assert.Nil(t, err)
}

func TestCompileErrorWithUnknownVariables(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", types.NewNumber(nil), ast.NewVariable("y"))}),
	)
	assert.Error(t, err)
}

func TestCompilePanicWithUntypedGlobals(t *testing.T) {
	assert.Panics(t, func() {
		Compile(
			ast.NewModule(
				"",
				[]ast.Bind{ast.NewBind("x", types.NewUnknown(nil), ast.NewNumber(42))}),
		)
	})
}

func TestCompileToCoreWithEmptySource(t *testing.T) {
	m, err := compileToCore(ast.NewModule("", []ast.Bind{}))
	assert.Nil(t, err)

	assert.Equal(t, coreast.NewModule("", nil, []coreast.Bind{}), m)
}

func TestCompileToCoreWithVariableBinds(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", types.NewNumber(nil), ast.NewNumber(42))},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"x",
					coreast.NewLambda(nil, true, nil, coreast.NewFloat64(42), coretypes.NewFloat64()),
				),
			},
		),
		m,
	)
}

func TestCompileToCoreWithFunctionBinds(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(
						types.NewNumber(nil),
						types.NewFunction(
							types.NewNumber(nil),
							types.NewNumber(nil),
							nil,
						),
						nil,
					),
					ast.NewLambda([]string{"x", "y"}, ast.NewNumber(42)),
				),
			},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"foo",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"foo.literal-0",
					coreast.NewLambda(nil, true, nil, coreast.NewFloat64(42), coretypes.NewFloat64()),
				),
				coreast.NewBind(
					"f",
					coreast.NewLambda(
						nil,
						false,
						[]coreast.Argument{
							coreast.NewArgument("x", coretypes.NewBoxed(coretypes.NewFloat64())),
							coreast.NewArgument("y", coretypes.NewBoxed(coretypes.NewFloat64())),
						},
						coreast.NewApplication(coreast.NewVariable("foo.literal-0"), nil),
						coretypes.NewBoxed(coretypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}

func TestCompileToCoreWithLetExpressions(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewLet(
						[]ast.Bind{ast.NewBind("y", types.NewNumber(nil), ast.NewNumber(42))},
						ast.NewVariable("y"),
					),
				),
			},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"foo",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"foo.literal-0",
					coreast.NewLambda(nil, true, nil, coreast.NewFloat64(42), coretypes.NewFloat64()),
				),
				coreast.NewBind(
					"x",
					coreast.NewLambda(
						nil,
						true,
						nil,
						coreast.NewLet(
							[]coreast.Bind{
								coreast.NewBind(
									"y",
									coreast.NewLambda(
										nil,
										true,
										nil,
										coreast.NewApplication(coreast.NewVariable("foo.literal-0"), nil),
										coretypes.NewBoxed(coretypes.NewFloat64()),
									),
								),
							},
							coreast.NewApplication(coreast.NewVariable("y"), nil),
						),
						coretypes.NewBoxed(coretypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}

func TestCompileToCoreWithLetExpressionsAndFreeVariables(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda(
						[]string{"x"},
						ast.NewLet(
							[]ast.Bind{ast.NewBind("y", types.NewNumber(nil), ast.NewVariable("x"))},
							ast.NewVariable("y"),
						),
					),
				),
			},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"foo",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"f",
					coreast.NewLambda(
						nil,
						false,
						[]coreast.Argument{
							coreast.NewArgument("x", coretypes.NewBoxed(coretypes.NewFloat64())),
						},
						coreast.NewLet(
							[]coreast.Bind{
								coreast.NewBind(
									"y",
									coreast.NewLambda(
										[]coreast.Argument{
											coreast.NewArgument("x", coretypes.NewBoxed(coretypes.NewFloat64())),
										},
										true,
										nil,
										coreast.NewApplication(coreast.NewVariable("x"), nil),
										coretypes.NewBoxed(coretypes.NewFloat64()),
									),
								),
							},
							coreast.NewApplication(coreast.NewVariable("y"), nil),
						),
						coretypes.NewBoxed(coretypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}

func TestCompileToCoreWithNestedLetExpressionsInLambdaExpressions(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"f",
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLambda(
						[]string{"x"},
						ast.NewLet(
							[]ast.Bind{
								ast.NewBind(
									"y", types.NewUnknown(nil), ast.NewLet(
										[]ast.Bind{ast.NewBind("z", types.NewUnknown(nil), ast.NewVariable("x"))},
										ast.NewVariable("z"),
									),
								),
							},
							ast.NewVariable("y"),
						),
					),
				),
			},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"foo",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"f",
					coreast.NewLambda(
						nil,
						false,
						[]coreast.Argument{
							coreast.NewArgument("x", coretypes.NewBoxed(coretypes.NewFloat64())),
						},
						coreast.NewLet(
							[]coreast.Bind{
								coreast.NewBind(
									"y",
									coreast.NewLambda(
										[]coreast.Argument{
											coreast.NewArgument("x", coretypes.NewBoxed(coretypes.NewFloat64())),
										},
										true,
										nil,
										coreast.NewLet(
											[]coreast.Bind{
												coreast.NewBind(
													"z",
													coreast.NewLambda(
														[]coreast.Argument{
															coreast.NewArgument(
																"x",
																coretypes.NewBoxed(coretypes.NewFloat64()),
															),
														},
														true,
														nil,
														coreast.NewApplication(coreast.NewVariable("x"), nil),
														coretypes.NewBoxed(coretypes.NewFloat64()),
													),
												),
											},
											coreast.NewApplication(coreast.NewVariable("z"), nil),
										),
										coretypes.NewBoxed(coretypes.NewFloat64()),
									),
								),
							},
							coreast.NewApplication(coreast.NewVariable("y"), nil),
						),
						coretypes.NewBoxed(coretypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}

func TestCompileToCoreWithBinaryOperations(t *testing.T) {
	m, err := compileToCore(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewBinaryOperation(ast.Add, ast.NewNumber(1), ast.NewNumber(1)),
				),
			},
		),
	)
	assert.Nil(t, err)

	assert.Equal(
		t,
		coreast.NewModule(
			"foo",
			nil,
			[]coreast.Bind{
				coreast.NewBind(
					"foo.literal-0",
					coreast.NewLambda(nil, true, nil, coreast.NewFloat64(1), coretypes.NewFloat64()),
				),
				coreast.NewBind(
					"foo.literal-1",
					coreast.NewLambda(nil, true, nil, coreast.NewFloat64(1), coretypes.NewFloat64()),
				),
				coreast.NewBind(
					"x",
					coreast.NewLambda(
						nil,
						true,
						nil,
						coreast.NewLet(
							[]coreast.Bind{
								coreast.NewBind(
									"result",
									coreast.NewLambda(
										nil,
										true,
										nil,
										coreast.NewPrimitiveCase(
											coreast.NewApplication(coreast.NewVariable("foo.literal-0"), nil),
											coretypes.NewBoxed(coretypes.NewFloat64()),
											nil,
											coreast.NewDefaultAlternative(
												"lhs",
												coreast.NewPrimitiveCase(
													coreast.NewApplication(coreast.NewVariable("foo.literal-1"), nil),
													coretypes.NewBoxed(coretypes.NewFloat64()),
													nil,
													coreast.NewDefaultAlternative(
														"rhs",
														coreast.NewPrimitiveOperation(
															coreast.AddFloat64,
															[]coreast.Atom{
																coreast.NewVariable("lhs"),
																coreast.NewVariable("rhs"),
															},
														),
													),
												),
											),
										),
										coretypes.NewFloat64(),
									),
								),
							},
							coreast.NewApplication(coreast.NewVariable("result"), nil),
						),
						coretypes.NewBoxed(coretypes.NewFloat64()),
					),
				),
			},
		),
		m,
	)
}

func TestCompileWithComplexBinaryOperations(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"foo",
			[]ast.Bind{
				ast.NewBind(
					"x",
					types.NewNumber(nil),
					ast.NewBinaryOperation(
						ast.Add,
						ast.NewNumber(1),
						ast.NewBinaryOperation(
							ast.Multiply,
							ast.NewNumber(2),
							ast.NewNumber(3),
						),
					),
				),
			},
		),
	)
	assert.Nil(t, err)
}
