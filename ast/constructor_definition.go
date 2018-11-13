package ast

// ConstructorDefinition is a constructor definition.
type ConstructorDefinition struct {
	name string
	tag  int
}

// NewConstructorDefinition creates a constructor definition.
func NewConstructorDefinition(n string, t int) ConstructorDefinition {
	return ConstructorDefinition{n, t}
}

// Name returns a name.
func (d ConstructorDefinition) Name() string {
	return d.name
}

// Tag returns a tag.
func (d ConstructorDefinition) Tag() int {
	return d.tag
}
