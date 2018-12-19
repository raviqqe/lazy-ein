package tinfer

import "github.com/ein-lang/ein/command/ast"

// InferTypes infers types in a module.
func InferTypes(m ast.Module) (ast.Module, error) {
	return newInferrer(m).Infer(m)
}
