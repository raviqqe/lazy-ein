package tinfer

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/debug"
	"github.com/ein-lang/ein/command/types"
)

// InferTypes infers types in a module.
func InferTypes(m ast.Module) (ast.Module, error) {
	m, err := newInferrer(m).Infer(m)

	if err != nil {
		return ast.Module{}, err
	} else if err := validateFullyTypedModule(m); err != nil {
		return ast.Module{}, err
	}

	return m, nil
}

func validateFullyTypedModule(m ast.Module) error {
	return m.VisitTypes(func(t types.Type) error {
		if v, ok := t.(types.Variable); ok {
			return debug.NewError("TypeInferenceError", "failed to infer a type", v.DebugInformation())
		}

		return nil
	})
}
