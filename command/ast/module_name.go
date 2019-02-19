package ast

import (
	"path/filepath"
	"strings"
)

// ModuleName is a unique module name.
type ModuleName (string)

// NewModuleName returns a module name.
func NewModuleName(f, rootDir string) (ModuleName, error) {
	d, err := filepath.Abs(rootDir)

	if err != nil {
		return "", err
	}

	f, err = filepath.Abs(f)

	if err != nil {
		return "", err
	}

	if !strings.HasPrefix(f, d) {
		return ModuleName(filepath.ToSlash(f)), nil
	}

	p, err := filepath.Rel(d, f)

	if err != nil {
		return "", err
	}

	return ModuleName(filepath.ToSlash(p)), nil
}