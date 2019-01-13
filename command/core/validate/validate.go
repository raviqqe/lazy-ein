package validate

import (
	"errors"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/validate/tcheck"
)

// Validate validates a module.
func Validate(m ast.Module) error {
	for _, f := range []func(ast.Module) error{checkFreeVariables, tcheck.CheckTypes} {
		if err := f(m); err != nil {
			return err
		}
	}

	// TODO: Check duplicate top-level names including constructors.

	return nil
}

func checkFreeVariables(m ast.Module) error {
	for _, b := range m.Binds() {
		if len(b.Lambda().FreeVariableNames()) != 0 {
			return errors.New("globals must not have free variables")
		}
	}

	return nil
}
