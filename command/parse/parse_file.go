package parse

import (
	"io/ioutil"

	"github.com/ein-lang/ein/command/ast"
)

// ParseFile parses a module file.
func ParseFile(f, rootDir string) (ast.Module, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return ast.Module{}, err
	}

	n, err := ast.NewModuleName(f, rootDir)

	if err != nil {
		return ast.Module{}, err
	}

	return Parse(string(bs), n)
}
