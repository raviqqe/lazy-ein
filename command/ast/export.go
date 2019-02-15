package ast

// Export is an export.
type Export struct {
	names []string
}

// NewExport creates an export.
func NewExport(ss ...string) Export {
	return Export{ss}
}

// Names returns names.
func (e Export) Names() []string {
	return e.names
}
