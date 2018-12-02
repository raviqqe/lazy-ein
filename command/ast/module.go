package ast

// Module is a module.
type Module struct {
	binds []Bind
}

// NewModule creates a module.
func NewModule(bs []Bind) Module {
	return Module{bs}
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}
