package generate

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestObjectCacheStore(t *testing.T) {
	d, err := ioutil.TempDir("", "")
	defer os.Remove(d)
	assert.Nil(t, err)

	f, err := ioutil.TempFile("", "")
	defer os.Remove(f.Name())
	assert.Nil(t, err)

	_, err = f.Write([]byte("bar"))
	assert.Nil(t, err)

	c := newObjectCache(d)

	s, ok, err := c.Get(f.Name())

	assert.Nil(t, err)
	assert.False(t, ok)
	assert.Equal(t, "", s)

	s, err = c.Store(f.Name(), []byte("baz"))

	assert.Nil(t, err)
	assert.NotEqual(t, "", s)

	ss, ok, err := c.Get(f.Name())

	assert.Nil(t, err)
	assert.True(t, ok)
	assert.Equal(t, s, ss)
}
