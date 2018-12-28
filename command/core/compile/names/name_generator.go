package names

import "fmt"

// NameGenerator is a name generator.
type NameGenerator struct {
	prefix  string
	indexes map[string]int
}

// NewNameGenerator creates a new name generator.
func NewNameGenerator(s string) *NameGenerator {
	return &NameGenerator{s, map[string]int{}}
}

// Generate generates a new name which is not duplicate.
func (g *NameGenerator) Generate(s string) string {
	s = g.prefix + "." + s

	if i, ok := g.indexes[s]; ok {
		g.indexes[s]++
		return s + "-" + fmt.Sprint(i+1)
	}

	g.indexes[s] = 0
	return s
}
