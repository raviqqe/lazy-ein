package names

import "strconv"

// NameGenerator is a name generator.
type NameGenerator struct {
	prefix string
	index  int
}

// NewNameGenerator creates a new name generator.
func NewNameGenerator(s string) *NameGenerator {
	return &NameGenerator{s, -1}
}

// Generate generates a new name which is not duplicate.
func (g *NameGenerator) Generate() string {
	g.index++
	return g.prefix + "-" + strconv.Itoa(g.index)
}
