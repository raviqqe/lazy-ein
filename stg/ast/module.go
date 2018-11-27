package ast

// Module is a module.
type Module struct {
	name                   string
	constructorDefinitions []ConstructorDefinition
	binds                  []Bind
}

// NewModule creates a module.
func NewModule(n string, ds []ConstructorDefinition, bs []Bind) Module {
	return Module{n, ds, bs}
}

// Name returns a name.
func (m Module) Name() string {
	return m.name
}

// ConstructorDefinitions returns constructor definitions.
func (m Module) ConstructorDefinitions() []ConstructorDefinition {
	return m.constructorDefinitions
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}
