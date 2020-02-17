package tcheck

import "github.com/raviqqe/lazy-ein/command/core/ast"

// CheckTypes checks types.
func CheckTypes(m ast.Module) error {
	return newTypeChecker().Check(m)
}
