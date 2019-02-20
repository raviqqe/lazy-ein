package metadata

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
)

// Module is a module metadata.
type Module struct {
	name          ast.ModuleName
	exportedBinds map[string]types.Type
}

// NewModule returns a module metadata.
func NewModule(m ast.Module) Module {
	ts := make(map[string]types.Type, len(m.Export().Names()))

	for _, b := range m.ExportedBinds() {
		ts[b.Name()] = types.Box(b.Type())
	}

	return Module{m.Name(), ts}
}

// Name returns a name.
func (m Module) Name() ast.ModuleName {
	return m.name
}

// ExportedBinds returns exported binds.
func (m Module) ExportedBinds() map[string]types.Type {
	return m.exportedBinds
}
