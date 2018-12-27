package tinfer

import (
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
)

type inferrer struct {
	variables         map[string]types.Type
	typeVariableCount *int
}

func newInferrer(m ast.Module) inferrer {
	c := 0
	i := inferrer{map[string]types.Type{}, &c}
	return i.addVariablesFromBinds(m.Binds())
}

func (i inferrer) Infer(m ast.Module) (ast.Module, error) {
	m = i.insertTypeVariables(m)
	ss := map[int]types.Type{}

	for _, b := range m.Binds() {
		t, es, err := i.inferType(b.Expression())

		if err != nil {
			return ast.Module{}, err
		}

		ees, err := b.Type().Unify(t)

		if err != nil {
			return ast.Module{}, err
		}

		sss, err := i.createSubstitutions(append(es, ees...))

		if err != nil {
			return ast.Module{}, err
		}

		for k, v := range sss {
			ss[k] = v
		}
	}

	return i.substituteVariablesInModule(m, ss), nil
}

func (i inferrer) inferType(e ast.Expression) (types.Type, []types.Equation, error) {
	es := []types.Equation{}

	switch e := e.(type) {
	case ast.Application:
		ee := ast.Expression(ast.NewApplication(e.Function(), e.Arguments()[:len(e.Arguments())-1]))

		if len(e.Arguments()) == 1 {
			ee = e.Function()
		}

		t, es, err := i.inferType(ee)

		if err != nil {
			return nil, nil, err
		}

		f, ok := t.(types.Function)

		if !ok {
			f = types.NewFunction(i.createTypeVariable(), i.createTypeVariable(), nil)
			ees, err := t.Unify(f)

			if err != nil {
				return nil, nil, err
			}

			es = append(es, ees...)
		}

		a, ees, err := i.inferType(e.Arguments()[len(e.Arguments())-1])

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		ees, err = f.Argument().Unify(a)

		if err != nil {
			return nil, nil, err
		}

		return f.Result(), append(es, ees...), nil
	case ast.BinaryOperation:
		es := []types.Equation{}

		for _, e := range []ast.Expression{e.LHS(), e.RHS()} {
			l, ees, err := i.inferType(e)

			if err != nil {
				return nil, nil, err
			}

			es = append(append(es, ees...), types.NewEquation(l, types.NewNumber(nil)))
		}

		return types.NewNumber(nil), es, nil
	case ast.Case:
		t, es, err := i.inferType(e.Expression())

		if err != nil {
			return nil, nil, err
		}

		ees, err := e.Type().Unify(t)

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		ts := make([]types.Type, 0, len(e.Alternatives())+1)

		for _, a := range e.Alternatives() {
			t, ees, err := i.inferType(a.Expression())

			if err != nil {
				return nil, nil, err
			}

			ts = append(ts, t)
			es = append(es, ees...)
		}

		if d, ok := e.DefaultAlternative(); ok {
			t, ees, err := i.addVariables(map[string]types.Type{d.Variable(): e.Type()}).inferType(
				d.Expression(),
			)

			if err != nil {
				return nil, nil, err
			}

			ts = append(ts, t)
			es = append(es, ees...)
		}

		for _, t := range ts[1:] {
			ees, err := ts[0].Unify(t)

			if err != nil {
				return nil, nil, err
			}

			es = append(es, ees...)
		}

		return ts[0], es, nil
	case ast.Lambda:
		as := make(map[string]types.Type, len(e.Arguments()))

		for _, s := range e.Arguments() {
			as[s] = i.createTypeVariable()
		}

		t, es, err := i.addVariables(as).inferType(e.Expression())

		if err != nil {
			return nil, nil, err
		}

		for i := len(as) - 1; i >= 0; i-- {
			t = types.NewFunction(as[e.Arguments()[i]], t, nil)
		}

		return t, es, nil
	case ast.Let:
		i = i.addVariablesFromBinds(e.Binds())

		for _, b := range e.Binds() {
			t, ees, err := i.inferType(b.Expression())

			if err != nil {
				return nil, nil, err
			}

			es = append(es, ees...)

			ees, err = b.Type().Unify(t)

			if err != nil {
				return nil, nil, err
			}

			es = append(es, ees...)
		}

		t, ees, err := i.inferType(e.Expression())

		if err != nil {
			return nil, nil, err
		}

		return t, append(es, ees...), nil
	case ast.Number:
		return types.NewNumber(nil), nil, nil
	case ast.Unboxed:
		t, es, err := i.inferType(e.Content())

		if err != nil {
			return nil, nil, err
		}

		return types.NewUnboxed(t, nil), es, nil
	case ast.Variable:
		t, ok := i.variables[e.Name()]

		if !ok {
			return nil, nil, fmt.Errorf("variable '%s' not found", e.Name())
		}

		return t, nil, nil
	}

	panic("unreachable")
}

func (i inferrer) createSubstitutions(es []types.Equation) (map[int]types.Type, error) {
	ss := map[int]types.Type{}

	for len(es) != 0 {
		e := es[0]
		es = es[1:]

		if v, ok := e.Left().(types.Variable); ok {
			ss = i.substituteVariablesInSubstitutions(ss, v, e.Right())
			es = i.substituteVariablesInEquations(es, v, e.Right())
			ss[v.Identifier()] = e.Right()
			continue
		} else if v, ok := e.Right().(types.Variable); ok {
			ss = i.substituteVariablesInSubstitutions(ss, v, e.Left())
			es = i.substituteVariablesInEquations(es, v, e.Left())
			ss[v.Identifier()] = e.Left()
			continue
		}

		ees, err := e.Left().Unify(e.Right())

		if err != nil {
			return nil, err
		}

		es = append(es, ees...)
	}

	return ss, nil
}

func (inferrer) substituteVariablesInSubstitutions(
	ss map[int]types.Type,
	v types.Variable,
	t types.Type,
) map[int]types.Type {
	sss := make(map[int]types.Type, len(ss))

	for i, tt := range ss {
		sss[i] = tt.SubstituteVariable(v, t)
	}

	return sss
}

func (inferrer) substituteVariablesInEquations(
	es []types.Equation,
	v types.Variable,
	t types.Type,
) []types.Equation {
	ees := make([]types.Equation, 0, len(es))

	for _, e := range es {
		ees = append(
			ees,
			types.NewEquation(e.Left().SubstituteVariable(v, t), e.Right().SubstituteVariable(v, t)),
		)
	}

	return ees
}

func (i inferrer) substituteVariablesInModule(m ast.Module, ss map[int]types.Type) ast.Module {
	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		switch e := e.(type) {
		case ast.Case:
			a, ok := e.DefaultAlternative()

			if !ok {
				return ast.NewCaseWithoutDefault(
					e.Expression(),
					i.substituteVariable(e.Type(), ss),
					e.Alternatives())
			}

			return ast.NewCase(
				e.Expression(),
				i.substituteVariable(e.Type(), ss),
				e.Alternatives(),
				a,
			)
		case ast.Let:
			bs := make([]ast.Bind, 0, len(e.Binds()))

			for _, b := range e.Binds() {
				bs = append(
					bs,
					ast.NewBind(b.Name(), i.substituteVariable(b.Type(), ss), b.Expression()),
				)
			}

			return ast.NewLet(bs, e.Expression())
		}

		return e
	}).(ast.Module)
}

func (inferrer) substituteVariable(t types.Type, ss map[int]types.Type) types.Type {
	return ss[t.(types.Variable).Identifier()]
}

func (i inferrer) addVariablesFromBinds(bs []ast.Bind) inferrer {
	m := make(map[string]types.Type, len(bs))

	for _, b := range bs {
		m[b.Name()] = types.Box(b.Type())
	}

	return i.addVariables(m)
}

func (i inferrer) addVariables(vs map[string]types.Type) inferrer {
	m := make(map[string]types.Type, len(i.variables)+len(vs))

	for k, v := range i.variables {
		m[k] = v
	}

	for k, v := range vs {
		m[k] = v
	}

	return inferrer{m, i.typeVariableCount}
}

func (i inferrer) createTypeVariable() types.Variable {
	t := types.NewVariable(*i.typeVariableCount, nil)
	*i.typeVariableCount++
	return t
}

func (i inferrer) insertTypeVariables(m ast.Module) ast.Module {
	return m.ConvertExpressions(func(e ast.Expression) ast.Expression {
		switch e := e.(type) {
		case ast.Case:
			a, ok := e.DefaultAlternative()

			if !ok {
				return ast.NewCaseWithoutDefault(
					e.Expression(),
					i.createTypeVariable(),
					e.Alternatives(),
				)
			}

			return ast.NewCase(e.Expression(), i.createTypeVariable(), e.Alternatives(), a)
		case ast.Let:
			bs := make([]ast.Bind, 0, len(e.Binds()))

			for _, b := range e.Binds() {
				bs = append(bs, ast.NewBind(b.Name(), i.createTypeVariable(), b.Expression()))
			}

			return ast.NewLet(bs, e.Expression())
		}

		return e
	}).(ast.Module)
}
