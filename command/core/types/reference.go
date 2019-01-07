package types

// Reference is a reference to a type.
type Reference struct {
	name string
}

// NewReference creates a reference to a type.
func NewReference(n string) Reference {
	return Reference{n}
}

// Name returns a name.
func (n Reference) Name() string {
	return n.name
}

func (Reference) isType() {}
