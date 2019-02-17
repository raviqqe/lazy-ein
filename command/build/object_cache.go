package build

import (
	"crypto/sha256"
	"encoding/base32"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ein-lang/ein/command/parse"
)

type objectCache struct {
	cacheDirectory, moduleRootDirectory string
}

func newObjectCache(cacheDir, rootDir string) objectCache {
	return objectCache{cacheDir, rootDir}
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
	h := sha256.New()

	if err := c.generateModuleHash(h, f); err != nil {
		return "", err
	}

	return filepath.Join(
		c.cacheDirectory,
		base32.HexEncoding.WithPadding(base32.NoPadding).EncodeToString(h.Sum(nil)),
	), nil
}

func (c objectCache) generateModuleHash(h hash.Hash, f string) error {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return err
	}

	h.Write([]byte(f))
	h.Write(bs)

	if err := c.generateSubmodulesHash(h, f, string(bs)); err != nil {
		return err
	}

	return nil
}

func (c objectCache) generateSubmodulesHash(h hash.Hash, f, s string) error {
	m, err := parse.Parse(f, s, c.moduleRootDirectory)

	if err != nil {
		return err
	}

	for _, i := range m.Imports() {
		if err := c.generateModuleHash(h, i.Path()); err != nil {
			return err
		}
	}

	return nil
}
