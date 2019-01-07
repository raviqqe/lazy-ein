package types

import (
	"testing"

	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
)

func TestListToTypeDefinition(t *testing.T) {
	assert.Equal(
		t,
		coreast.NewTypeDefinition(
			"$List.$Number.$end",
			coretypes.NewAlgebraic(
				[]coretypes.Constructor{
					coretypes.NewConstructor(
						"$Cons.$Number.$end",
						[]coretypes.Type{
							coretypes.NewBoxed(coretypes.NewFloat64()),
							coretypes.NewBoxed(coretypes.NewReference("$List.$Number.$end")),
						},
					),
					coretypes.NewConstructor("$Nil.$Number.$end", nil),
				},
			),
		),
		NewList(NewNumber(nil), nil).ToTypeDefinition(),
	)
}