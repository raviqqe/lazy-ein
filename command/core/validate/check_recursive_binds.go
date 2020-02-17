package validate

import (
	"errors"

	"github.com/raviqqe/lazy-ein/command/core/ast"
)

func checkRecursiveBinds(m ast.Module) error {
	for _, b := range m.Binds() {
		if err := checkRecursiveBind(b); err != nil {
			return err
		}
	}

	return m.VisitExpressions(
		func(e ast.Expression) error {
			switch e := e.(type) {
			case ast.Let:
				for _, b := range e.Binds() {
					if err := checkRecursiveBind(b); err != nil {
						return err
					}
				}
			}

			return nil
		},
	)
}

func checkRecursiveBind(b ast.Bind) error {
	if a, ok := b.Lambda().Body().(ast.FunctionApplication); ok && b.Name() == a.Function().Name() {
		return errors.New("invalid recursive bind")
	}

	return nil
}
