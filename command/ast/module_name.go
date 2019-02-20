package ast

import (
	"path"
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

// NewModuleNameFromString returns a module name of a string.
func NewModuleNameFromString(s string) ModuleName {
	return ModuleName(s)
}

// Qualify qualifies a name in a module.
func (n ModuleName) Qualify(s string) string {
	return path.Base(string(n)) + "." + s
}

// FullyQualify qualifies a name in a module.
func (n ModuleName) FullyQualify(s string) string {
	return string(n) + "." + s
}

// ToPath converts a module name to a path.
func (n ModuleName) ToPath(rootDir string) string {
	s := string(n)

	if path.IsAbs(s) {
		return s
	}

	return path.Join(rootDir, s)
}
