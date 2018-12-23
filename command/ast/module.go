package ast

// Module is a module.
type Module struct {
	name  string
	binds []Bind
}

// NewModule creates a module.
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

// ConvertExpressions visits expressions.
func (m Module) ConvertExpressions(f func(Expression) Expression) node {
	bs := make([]Bind, 0, len(m.binds))

	for _, b := range m.binds {
		bs = append(bs, b.ConvertExpressions(f).(Bind))
	}

	return NewModule(m.name, bs)
}
