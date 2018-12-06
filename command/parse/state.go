package parse

import (
	"strings"

	"github.com/raviqqe/jsonxx/command/debug"
	"github.com/raviqqe/parcom"
)

type state struct {
	*parcom.PositionalState
	filename, source string
}

func newState(f, s string) *state {
	return &state{parcom.NewPositionalState(s), f, s}
}

func (s state) debugInformation() *debug.Information {
	return debug.NewInformation(
		s.filename,
		s.Line(),
		s.Column(),
		strings.Split(s.source, "\n")[s.Line()-1],
	)
}