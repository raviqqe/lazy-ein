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
					[]string{"x"},
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewVariable("x"),
				),
				ast.NewBind(
					"x",
					nil,
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

func TestCompileErrorWithUnknownVariables(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", nil, types.NewNumber(nil), ast.NewVariable("y"))}),
	)
	assert.Error(t, err)
}

func TestCompileErrorWithUntypedGlobals(t *testing.T) {
	_, err := Compile(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", nil, types.NewVariable(nil), ast.NewNumber(42))}),
	)
	assert.Error(t, err)
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
			[]ast.Bind{ast.NewBind("x", nil, types.NewNumber(nil), ast.NewNumber(42))},
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
					[]string{"x", "y"},
					types.NewFunction(
						types.NewNumber(nil),
						types.NewFunction(
							types.NewNumber(nil),
							types.NewNumber(nil),
							nil,
						),
						nil,
					),
					ast.NewNumber(42),
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
					nil,
					types.NewNumber(nil),
					ast.NewLet(
						[]ast.Bind{ast.NewBind("y", nil, types.NewNumber(nil), ast.NewNumber(42))},
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
					[]string{"x"},
					types.NewFunction(types.NewNumber(nil), types.NewNumber(nil), nil),
					ast.NewLet(
						[]ast.Bind{ast.NewBind("y", nil, types.NewNumber(nil), ast.NewVariable("x"))},
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
