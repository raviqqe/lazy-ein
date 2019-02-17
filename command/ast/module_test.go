package ast

import (
	"testing"

	"github.com/ein-lang/ein/command/types"
	"github.com/stretchr/testify/assert"
)

func TestModuleIsMainModule(t *testing.T) {
	assert.True(
		t,
		NewModule(
			"",
			NewExport(),
			nil,
			[]Bind{NewBind("main", types.NewNumber(nil), NewNumber(42))},
		).IsMainModule(),
	)
	assert.False(
		t,
		NewModule(
			"",
			NewExport(),
			nil,
			[]Bind{NewBind("x", types.NewNumber(nil), NewNumber(42))},
		).IsMainModule(),
	)
}
