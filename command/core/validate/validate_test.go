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
	assert.Nil(t, validate.Validate(ast.NewModule("foo", nil)))
}

func TestValidateError(t *testing.T) {
	tt := types.NewAlgebraic(types.NewConstructor())

	assert.Equal(
		t,
		errors.New("globals must not have free variables"),
		validate.Validate(
			ast.NewModule(
				"foo",
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
							true,
							ast.NewConstructorApplication(ast.NewConstructor(tt, 0), nil),
							tt,
						),
					),
				},
			),
		),
	)
}
