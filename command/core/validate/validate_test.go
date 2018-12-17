package validate_test

import (
	"errors"
	"testing"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/core/validate"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert.Nil(t, validate.Validate(ast.NewModule("foo", nil, nil)))
}

func TestValidateError(t *testing.T) {
	assert.Equal(
		t,
		errors.New("globals must not have free variables"),
		validate.Validate(
			ast.NewModule(
				"foo",
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewLambda(
							[]ast.Argument{ast.NewArgument("x", types.NewBoxed(types.NewFloat64()))},
							true,
							nil,
							ast.NewApplication(ast.NewVariable("x"), nil),
							types.NewBoxed(types.NewFloat64()),
						),
					),
				},
			),
		),
	)
}
