package validate

import (
	"errors"

	"github.com/ein-lang/ein/command/core/ast"
)

func checkFreeVariables(m ast.Module) error {
	for _, b := range m.Binds() {
		if len(b.Lambda().FreeVariableNames()) != 0 {
			return errors.New("globals must not have free variables")
		}
	}

	return nil
}
