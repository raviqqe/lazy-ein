package types

import "fmt"

// Boxed is a boxed type.
type Boxed struct {
	content Type
}

// NewBoxed creates a boxed type.
func NewBoxed(t Type) Boxed {
	if _, ok := t.(Boxed); ok {
		panic("cannot box boxed types")
	}

	return Boxed{t}
}

// Content returns a content type.
func (b Boxed) Content() Type {
	return b.content
}

func (b Boxed) String() string {
	return fmt.Sprintf("Boxed(%v)", b.content)
}

func (b Boxed) equal(t Type) bool {
	bb, ok := t.(Boxed)

	if !ok {
		return false
	}

	return b.content.equal(bb.content)
}
