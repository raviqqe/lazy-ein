package generate_test

import (
	"os"
	"testing"

	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/types"
	"github.com/raviqqe/jsonxx/command/generate"
	"github.com/stretchr/testify/assert"
)

func TestExecutable(t *testing.T) {
	os.Setenv("JSONXX_ROOT", "../..")

	assert.Nil(
		t,
		generate.Executable(
			"main.jsonxx",
			ast.NewModule(
				"main.jsonxx",
				nil,
				[]ast.Bind{
					ast.NewBind(
						"main",
						ast.NewLambda(
							nil,
							false,
							[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
							ast.NewApplication(ast.NewVariable("x"), nil),
							types.NewBoxed(types.NewFloat64()),
						),
					),
				},
			),
		),
	)
}

func TestExecutableErrorWithoutMainFunction(t *testing.T) {
	os.Setenv("JSONXX_ROOT", "../..")

	assert.Error(
		t,
		generate.Executable("main.jsonxx", ast.NewModule("main.jsonxx", nil, nil)),
	)
}

func TestExecutableErrorWithoutRootEnvironmentVariable(t *testing.T) {
	os.Unsetenv("JSONXX_ROOT")

	assert.Error(
		t,
		generate.Executable(
			"main.jsonxx",
			ast.NewModule(
				"main.jsonxx",
				nil,
				[]ast.Bind{
					ast.NewBind(
						"main",
						ast.NewLambda(
							nil,
							false,
							[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
							ast.NewApplication(ast.NewVariable("x"), nil),
							types.NewBoxed(types.NewFloat64()),
						),
					),
				},
			),
		),
	)
}
