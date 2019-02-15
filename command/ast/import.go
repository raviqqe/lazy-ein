package ast

// Import is a module import.
type Import struct {
	path string
}

// NewImport creates a module import.
func NewImport(s string) Import {
	return Import{s}
}

// Path returns a path.
func (i Import) Path() string {
	return i.path
}
