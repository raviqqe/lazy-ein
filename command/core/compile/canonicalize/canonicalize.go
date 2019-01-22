package canonicalize

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
)

// Canonicalize canonicalizes a module.
func Canonicalize(m ast.Module) ast.Module {
	ts := []types.Type{}

	m.ConvertTypes(func(t types.Type) types.Type {
		switch t.(type) {
		case types.Algebraic, types.Function:
			if newTypeValidator().Validate(t) {
				ts = append(ts, newTypeCanonicalizer().Canonicalize(t))
			}
		}

		return t
	})

	return m.ConvertTypes(func(t types.Type) types.Type {
		if !newTypeValidator().Validate(t) {
			return t
		}

		for _, tt := range ts {
			if types.Equal(t, tt) {
				return tt
			}
		}

		return t
	})
}
