package ast

// Import is a module import.
type Import struct {
	name ModuleName
}

// NewImport creates a module import.
func NewImport(n ModuleName) Import {
	return Import{n}
}

// Name returns a name.
func (i Import) Name() ModuleName {
	return i.name
}
