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
	m := compileToCore(ast.NewModule("", []ast.Bind{}))

	assert.Equal(t, coreast.NewModule("", nil, []coreast.Bind{}), m)
}

func TestCompileWithVariableBinds(t *testing.T) {
	m := compileToCore(
		ast.NewModule(
			"",
			[]ast.Bind{ast.NewBind("x", nil, types.NewNumber(nil), ast.NewNumber(42))},
		),
	)

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

func TestCompileWithFunctionBinds(t *testing.T) {
	m := compileToCore(
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

func TestCompileWithLetExpressions(t *testing.T) {
	m := compileToCore(
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
