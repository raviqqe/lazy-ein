package types

// Boxed is a boxed type.
type Boxed struct {
	content Type
}

// NewBoxed creates a new boxed type.
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

func (Boxed) isType() {}
