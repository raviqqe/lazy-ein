package names

import "fmt"

// NameGenerator is a name generator.
type NameGenerator struct {
	prefix  string
	indexes map[string]int
}

// NewNameGenerator creates a new name generator.
func NewNameGenerator(s string) NameGenerator {
	if s != "" {
		s += "."
	}

	return NameGenerator{s, map[string]int{}}
}

// Generate generates a new name which is not duplicate.
func (g NameGenerator) Generate(s string) string {
	i := g.indexes[s]
	g.indexes[s]++
	return fmt.Sprintf("%v%v-%v", g.prefix, s, i)
}
