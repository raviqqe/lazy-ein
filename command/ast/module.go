package ast

import "github.com/ein-lang/ein/command/types"

// Module is a module.
type Module struct {
	name    ModuleName
	export  Export
	imports []Import
	binds   []Bind
}

// NewModule creates a module.
func NewModule(n ModuleName, e Export, is []Import, bs []Bind) Module {
	return Module{n, e, is, bs}
}

// Name returns a name.
func (m Module) Name() ModuleName {
	return m.name
}

// Export returns an export.
func (m Module) Export() Export {
	return m.export
}

// Imports returns imports.
func (m Module) Imports() []Import {
	return m.imports
}

// Binds returns binds.
func (m Module) Binds() []Bind {
	return m.binds
}

// IsMainModule checks if it is a main module.
func (m Module) IsMainModule() bool {
	for _, b := range m.binds {
		if b.Name() == MainFunctionName {
			return true
		}
	}

	return false
}

// ConvertExpressions converts expressions.
func (m Module) ConvertExpressions(f func(Expression) Expression) Node {
	bs := make([]Bind, 0, len(m.binds))

	for _, b := range m.binds {
		bs = append(bs, b.ConvertExpressions(f).(Bind))
	}

	return NewModule(m.name, m.export, m.imports, bs)
}

// VisitTypes visits types.
func (m Module) VisitTypes(f func(types.Type) error) error {
	for _, b := range m.binds {
		if err := b.VisitTypes(f); err != nil {
			return err
		}
	}

	return nil
}
