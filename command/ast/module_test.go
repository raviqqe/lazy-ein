package ast

import (
	"testing"

	"github.com/raviqqe/lazy-ein/command/types"
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

func TestModuleExportedBinds(t *testing.T) {
	assert.Equal(
		t,
		[]Bind{NewBind("x", types.NewNumber(nil), NewNumber(42))},
		NewModule(
			"",
			NewExport("x"),
			nil,
			[]Bind{NewBind("x", types.NewNumber(nil), NewNumber(42))},
		).ExportedBinds(),
	)
	assert.Equal(
		t,
		[]Bind{},
		NewModule(
			"",
			NewExport(),
			nil,
			[]Bind{NewBind("x", types.NewNumber(nil), NewNumber(42))},
		).ExportedBinds(),
	)
}
