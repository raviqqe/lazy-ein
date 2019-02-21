package build

import (
	"crypto/sha256"
	"encoding/base32"
	"hash"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/parse"
	"llvm.org/llvm/bindings/go/llvm"
)

type objectCache struct {
	cacheDirectory, moduleRootDirectory string
}

func newObjectCache(cacheDir, rootDir string) objectCache {
	return objectCache{cacheDir, rootDir}
}

func (c objectCache) Store(fname string, m llvm.Module) error {
	p, err := c.generatePath(fname)

	if err != nil {
		return err
	}

	f, err := os.Create(p)

	if err != nil {
		return err
	}

	return llvm.WriteBitcodeToFile(m, f)
}

func (c objectCache) Get(fname string) (llvm.Module, bool, error) {
	p, err := c.generatePath(fname)

	if err != nil {
		return llvm.Module{}, false, err
	}

	if _, err := os.Stat(p); err != nil {
		return llvm.Module{}, false, nil
	}

	m, err := llvm.ParseBitcodeFile(p)

	if err != nil {
		return llvm.Module{}, false, err
	}

	return m, true, nil
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

	m, err := parse.Parse(f, c.moduleRootDirectory)

	if err != nil {
		return err
	}

	h.Write([]byte(m.Name()))
	h.Write(bs)

	if err := c.generateSubmodulesHash(h, m); err != nil {
		return err
	}

	return nil
}

func (c objectCache) generateSubmodulesHash(h hash.Hash, m ast.Module) error {
	for _, i := range m.Imports() {
		if err := c.generateModuleHash(h, i.Name().ToPath(c.moduleRootDirectory)); err != nil {
			return err
		}
	}

	return nil
}
