package validate

import (
	"github.com/raviqqe/lazy-ein/command/core/ast"
	"github.com/raviqqe/lazy-ein/command/core/validate/tcheck"
)

// Validate validates a module.
func Validate(m ast.Module) error {
	for _, f := range []func(ast.Module) error{
		checkFreeVariables,
		checkRecursiveBinds,
		tcheck.CheckTypes,
	} {
		if err := f(m); err != nil {
			return err
		}
	}

	// TODO: Check duplicate top-level names including constructors.

	return nil
}
