package compile

import "github.com/ein-lang/ein/command/ast"

type freeVariableFinder struct {
	variables map[string]struct{}
}

func newFreeVariableFinder(vs map[string]struct{}) freeVariableFinder {
	return freeVariableFinder{vs}
}

func (f freeVariableFinder) Find(e ast.Expression) []string {
	switch e := e.(type) {
	case ast.Application:
		ss := f.Find(e.Function())

		for _, a := range e.Arguments() {
			ss = append(ss, f.Find(a)...)
		}

		return ss
	case ast.Lambda:
		ss := make([]string, 0, len(e.Arguments()))

		for _, s := range e.Arguments() {
			ss = append(ss, s)
		}

		return f.addVariables(ss...).Find(e.Expression())
	case ast.Let:
		ss := make([]string, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			ss = append(ss, b.Name())
		}

		return f.addVariables(ss...).Find(e.Expression())
	case ast.Unboxed:
		return nil
	case ast.Variable:
		if _, ok := f.variables[e.Name()]; ok {
			return nil
		}

		return []string{e.Name()}
	}

	panic("unreahable")
}

func (f freeVariableFinder) addVariables(ss ...string) freeVariableFinder {
	m := make(map[string]struct{}, len(f.variables)+len(ss))

	for k := range m {
		m[k] = struct{}{}
	}

	for _, s := range ss {
		m[s] = struct{}{}
	}

	return freeVariableFinder{m}
}
