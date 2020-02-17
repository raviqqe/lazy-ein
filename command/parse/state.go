package parse

import (
	"strings"

	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/debug"
	"github.com/raviqqe/parcom"
)

type state struct {
	*parcom.PositionalState
	source     string
	moduleName ast.ModuleName
}

func newState(s string, n ast.ModuleName) *state {
	return &state{parcom.NewPositionalState(s), s, n}
}

func (s state) debugInformation() *debug.Information {
	return debug.NewInformation(
		string(s.moduleName),
		s.Line(),
		s.Column(),
		strings.Split(s.source, "\n")[s.Line()-1],
	)
}
