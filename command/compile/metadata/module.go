package metadata

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	coreast "github.com/raviqqe/lazy-ein/command/core/ast"
	coretypes "github.com/raviqqe/lazy-ein/command/core/types"
	"github.com/raviqqe/lazy-ein/command/types"
)

// Module is a module metadata.
type Module struct {
	name          ast.ModuleName
	exportedBinds map[string]types.Type
	declarations  []coreast.Declaration
}

// NewModule returns a module metadata.
func NewModule(m ast.Module) Module {
	ts := make(map[string]types.Type, len(m.Export().Names()))

	for _, b := range m.ExportedBinds() {
		ts[b.Name()] = types.Box(b.Type())
	}

	ds := make([]coreast.Declaration, 0, len(ts))

	for n, t := range ts {
		n := m.Name().Qualify(n)

		switch t := t.(type) {
		case types.Function:
			f := t.ToCore().(coretypes.Function)

			ds = append(
				ds,
				coreast.NewDeclaration(
					n,
					coreast.NewLambdaDeclaration(nil, f.Arguments(), f.Result()),
				),
			)
		default:
			ds = append(
				ds,
				coreast.NewDeclaration(
					n,
					coreast.NewLambdaDeclaration(nil, nil, coretypes.Unbox(t.ToCore())),
				),
			)
		}
	}

	return Module{m.Name(), ts, ds}
}

// Name returns a name.
func (m Module) Name() ast.ModuleName {
	return m.name
}

// ExportedBinds returns exported binds.
func (m Module) ExportedBinds() map[string]types.Type {
	return m.exportedBinds
}

// Declarations reutrns declarations.
func (m Module) Declarations() []coreast.Declaration {
	return m.declarations
}
