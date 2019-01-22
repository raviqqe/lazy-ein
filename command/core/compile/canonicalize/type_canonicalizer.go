package canonicalize

import "github.com/ein-lang/ein/command/core/types"

type typeCanonicalizer struct {
	environment []types.Type
}

func newTypeCanonicalizer() typeCanonicalizer {
	return typeCanonicalizer{}
}

func (c typeCanonicalizer) Canonicalize(t types.Type) types.Type {
	switch t := t.(type) {
	case types.Algebraic:
		c = c.pushType(t)
		cs := make([]types.Constructor, 0, len(t.Constructors()))

		for _, cc := range t.Constructors() {
			es := []types.Type(nil) // Normalize empty slices for testing.

			for _, e := range cc.Elements() {
				es = append(es, c.canonicalizeInner(e))
			}

			cs = append(cs, types.NewConstructor(es...))
		}

		return types.NewAlgebraic(cs[0], cs[1:]...)
	case types.Boxed:
		return types.NewBoxed(c.canonicalizeInner(t.Content()).(types.Boxable))
	case types.Function:
		c = c.pushType(t)
		as := make([]types.Type, 0, len(t.Arguments()))

		for _, a := range t.Arguments() {
			as = append(as, c.canonicalizeInner(a))
		}

		return types.NewFunction(as, c.Canonicalize(t.Result()))
	}

	return t
}

func (c typeCanonicalizer) canonicalizeInner(t types.Type) types.Type {
	for i, tt := range c.environment {
		if types.EqualWithEnvironment(t, tt, c.environment) {
			return types.NewIndex(i)
		}
	}

	return c.Canonicalize(t)
}

func (c typeCanonicalizer) pushType(t types.Type) typeCanonicalizer {
	return typeCanonicalizer{append(c.environment, t)}
}
