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
	case ast.Let:
		ss := make([]string, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			ss = append(ss, b.Name())
		}

		return f.addLocalVariables(ss...).Find(e.Expression())
	case ast.Number:
		return nil
	case ast.Variable:
		if _, ok := f.variables[e.Name()]; ok {
			return nil
		}

		return []string{e.Name()}
	}

	panic("unreahable")
}

func (f freeVariableFinder) addLocalVariables(ss ...string) freeVariableFinder {
	m := make(map[string]struct{}, len(f.variables)+len(ss))

	for k := range m {
		m[k] = struct{}{}
	}

	for _, s := range ss {
		m[s] = struct{}{}
	}

	return freeVariableFinder{m}
}
