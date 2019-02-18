package ast

import "path/filepath"

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

	p, err := filepath.Rel(d, f)

	if err != nil {
		return "", err
	}

	return ModuleName(filepath.ToSlash(p)), nil
}
