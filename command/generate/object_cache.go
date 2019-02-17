package generate

import (
	"crypto/sha256"
	"encoding/base32"
	"io/ioutil"
	"os"
	"path/filepath"
)

type objectCache struct {
	directory string
}

func newObjectCache(s string) objectCache {
	return objectCache{s}
}

func (c objectCache) Store(f string, bs []byte) (string, error) {
	p, err := c.generatePath(f)

	if err != nil {
		return "", err
	}

	return p, ioutil.WriteFile(p, bs, 0644)
}

func (c objectCache) Get(f string) (string, bool, error) {
	p, err := c.generatePath(f)

	if err != nil {
		return "", false, err
	}

	if _, err := os.Stat(p); err != nil {
		return "", false, nil
	}

	return p, true, nil
}

func (c objectCache) generatePath(f string) (string, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return "", err
	}

	f, err = filepath.Abs(f)

	if err != nil {
		return "", err
	}

	h := sha256.New()
	h.Write([]byte(f))
	h.Write(bs)

	return filepath.Join(
		c.directory,
		base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(h.Sum(nil)),
	), nil
}
