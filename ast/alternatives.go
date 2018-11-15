package ast

// Alternatives are alternatives.
type Alternatives interface {
	toSlice() []Alternative
}
