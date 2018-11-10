package ast

// Module is a module.
type Module struct {
	name  string
	binds []Bind
}

// NewModule creates a new module.
func NewModule(n string, bs []Bind) Module {
	return Module{n, bs}
}

// Name returns a name.
func (m Module) Name() string {
	return m.name
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}
