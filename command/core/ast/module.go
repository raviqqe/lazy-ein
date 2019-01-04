package ast

// Module is a module.
type Module struct {
	name            string
	typeDefinitions []TypeDefinition
	binds           []Bind
}

// NewModule creates a module.
func NewModule(n string, ds []TypeDefinition, bs []Bind) Module {
	return Module{n, ds, bs}
}

// Name returns a name.
func (m Module) Name() string {
	return m.name
}

// TypeDefinitions returns constructor definitions.
func (m Module) TypeDefinitions() []TypeDefinition {
	return m.typeDefinitions
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}
