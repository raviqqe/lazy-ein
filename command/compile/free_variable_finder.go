package compile

import "github.com/raviqqe/lazy-ein/command/ast"

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
	case ast.BinaryOperation:
		return append(f.Find(e.LHS()), f.Find(e.RHS())...)
	case ast.Case:
		ss := f.Find(e.Argument())

		for _, a := range e.Alternatives() {
			ss = append(ss, f.addVariablesFromPattern(a.Pattern()).Find(a.Expression())...)
		}

		if d, ok := e.DefaultAlternative(); ok {
			ss = append(ss, f.addVariables(d.Variable()).Find(d.Expression())...)
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

		f = f.addVariables(ss...)

		sss := []string{}

		for _, b := range e.Binds() {
			sss = append(sss, f.Find(b.Expression())...)
		}

		return append(sss, f.Find(e.Expression())...)
	case ast.List:
		ss := []string{}

		for _, a := range e.Arguments() {
			ss = append(ss, f.Find(a.Expression())...)
		}

		return ss
	case ast.Number:
		break
	case ast.Unboxed:
		return nil
	case ast.Variable:
		if _, ok := f.variables[e.Name()]; ok {
			return nil
		}

		return []string{e.Name()}
	}

	panic("unreachable")
}

func (f freeVariableFinder) addVariablesFromPattern(e ast.Expression) freeVariableFinder {
	l, ok := e.(ast.List)

	if !ok {
		return f
	}

	ss := make([]string, 0, len(l.Arguments()))

	for _, a := range l.Arguments() {
		if v, ok := a.Expression().(ast.Variable); ok {
			ss = append(ss, v.Name())
		}
	}

	return f.addVariables(ss...)
}

func (f freeVariableFinder) addVariables(ss ...string) freeVariableFinder {
	m := make(map[string]struct{}, len(f.variables)+len(ss))

	for k := range f.variables {
		m[k] = struct{}{}
	}

	for _, s := range ss {
		m[s] = struct{}{}
	}

	return freeVariableFinder{m}
}
