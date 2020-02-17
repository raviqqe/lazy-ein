package tinfer

import (
	"github.com/raviqqe/lazy-ein/command/ast"
	"github.com/raviqqe/lazy-ein/command/compile/metadata"
	"github.com/raviqqe/lazy-ein/command/debug"
	"github.com/raviqqe/lazy-ein/command/types"
)

// InferTypes infers types in a module with imported modules.
func InferTypes(m ast.Module, ms []metadata.Module) (ast.Module, error) {
	m, err := newInferrer(m, ms).Infer(m)

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
