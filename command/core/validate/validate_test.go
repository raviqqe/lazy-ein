package validate_test

import (
	"errors"
	"testing"

	"github.com/raviqqe/lazy-ein/command/core/ast"
	"github.com/raviqqe/lazy-ein/command/core/types"
	"github.com/raviqqe/lazy-ein/command/core/validate"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	assert.Nil(t, validate.Validate(ast.NewModule(nil, nil)))
}

func TestValidateError(t *testing.T) {
	tt := types.NewAlgebraic(types.NewConstructor())

	assert.Equal(
		t,
		errors.New("globals must not have free variables"),
		validate.Validate(
			ast.NewModule(
				nil,
				[]ast.Bind{
					ast.NewBind(
						"x",
						ast.NewVariableLambda(
							[]ast.Argument{ast.NewArgument("x", types.NewFloat64())},
							ast.NewConstructorApplication(ast.NewConstructor(tt, 0), nil),
							tt,
						),
					),
				},
			),
		),
	)
}
