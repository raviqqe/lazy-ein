package types

// Named is a reference to a named type.
type Named struct {
	name string
}

// NewNamed creates a named type.
func NewNamed(n string) Named {
	return Named{n}
}

// Name returns a name.
func (n Named) Name() string {
	return n.name
}

func (Named) isType() {}
