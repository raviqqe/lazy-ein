package parse

import "github.com/raviqqe/parcom"

type state struct {
	*parcom.State
	filename string
}

func newState(f, s string) *state {
	return &state{parcom.NewState(s), f}
}
