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
	} else if _, ok := t.(Function); ok {
		panic("cannot box function types")
	}

	return Boxed{t}
}

// Content returns a content type.
func (b Boxed) Content() Type {
	return b.content
}

// ConvertTypes converts types.
func (b Boxed) ConvertTypes(f func(Type) Type) Type {
	return f(Boxed{b.content.ConvertTypes(f)})
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
