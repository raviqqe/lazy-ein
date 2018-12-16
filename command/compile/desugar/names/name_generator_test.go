package names

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameGeneratorGenerate(t *testing.T) {
	g := NewNameGenerator("")

	assert.Equal(t, "foo-0", g.Generate("foo"))
	assert.Equal(t, "bar-0", g.Generate("bar"))
	assert.Equal(t, "foo-1", g.Generate("foo"))
}

func TestNameGeneratorGenerateWithPrefix(t *testing.T) {
	g := NewNameGenerator("my")

	assert.Equal(t, "my.foo-0", g.Generate("foo"))
}
