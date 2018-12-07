package names

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameGeneratorGenerate(t *testing.T) {
	g := NewNameGenerator("foo")

	assert.Equal(t, "foo-0", g.Generate())
	assert.Equal(t, "foo-1", g.Generate())
	assert.Equal(t, "foo-2", g.Generate())
}
