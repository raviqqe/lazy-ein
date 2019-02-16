package ast

import "github.com/ein-lang/ein/command/core/types"

// Module is a module.
type Module struct {
	name         string
	declarations []Bind
	binds        []Bind
}

// NewModule creates a module.
func NewModule(s string, ds []Bind, bs []Bind) Module {
	return Module{s, ds, bs}
}

// Name returns a name.
func (m Module) Name() string {
	return m.name
}

// Declarations returns declarations.
func (m Module) Declarations() []Bind {
	return m.declarations
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}

// Types returns types.
func (m Module) Types() []types.Type {
	ts := map[string]types.Type{}

	m.ConvertTypes(func(t types.Type) types.Type {
		ts[t.String()] = t
		return t
	})

	tts := make([]types.Type, 0, len(ts))

	for _, t := range ts {
		tts = append(tts, t)
	}

	return tts
}

// ConvertTypes converts types.
func (m Module) ConvertTypes(f func(types.Type) types.Type) Module {
	bs := make([]Bind, 0, len(m.binds))

	for _, b := range m.binds {
		bs = append(bs, b.ConvertTypes(f))
	}

	return Module{m.name, m.declarations, bs}
}
