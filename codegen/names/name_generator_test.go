package names

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNameGeneratorGenerate(t *testing.T) {
	g := NewNameGenerator("foo")

	assert.Equal(t, "foo.bar", g.Generate("bar"))
	assert.Equal(t, "foo.baz", g.Generate("baz"))
	assert.Equal(t, "foo.bar.1", g.Generate("bar"))
}
