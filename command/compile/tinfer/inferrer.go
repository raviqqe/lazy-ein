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
	switch e := e.(type) {
	case ast.Application:
		return i.inferApplication(e)
	case ast.BinaryOperation:
		return i.inferBinaryOperation(e)
	case ast.Case:
		return i.inferCase(e)
	case ast.Lambda:
		return i.inferLambda(e)
	case ast.Let:
		return i.inferLet(e)
	case ast.List:
		return i.inferList(e)
	case ast.Number:
		return types.NewNumber(nil), nil, nil
	case ast.Unboxed:
		return i.inferUnboxed(e)
	case ast.Variable:
		return i.inferVariable(e)
	}

	panic("unreachable")
}

func (i inferrer) inferApplication(a ast.Application) (types.Type, []types.Equation, error) {
	e := ast.Expression(ast.NewApplication(a.Function(), a.Arguments()[:len(a.Arguments())-1]))

	if len(a.Arguments()) == 1 {
		e = a.Function()
	}

	t, es, err := i.inferType(e)

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

	arg, ees, err := i.inferType(a.Arguments()[len(a.Arguments())-1])

	if err != nil {
		return nil, nil, err
	}

	es = append(es, ees...)

	ees, err = f.Argument().Unify(arg)

	if err != nil {
		return nil, nil, err
	}

	return f.Result(), append(es, ees...), nil
}

func (i inferrer) inferBinaryOperation(o ast.BinaryOperation) (types.Type, []types.Equation, error) {
	es := []types.Equation{}

	for _, e := range []ast.Expression{o.LHS(), o.RHS()} {
		l, ees, err := i.inferType(e)

		if err != nil {
			return nil, nil, err
		}

		es = append(append(es, ees...), types.NewEquation(l, types.NewNumber(nil)))
	}

	return types.NewNumber(nil), es, nil
}

func (i inferrer) inferCase(c ast.Case) (types.Type, []types.Equation, error) {
	t, es, err := i.inferType(c.Expression())

	if err != nil {
		return nil, nil, err
	}

	ees, err := c.Type().Unify(t)

	if err != nil {
		return nil, nil, err
	}

	es = append(es, ees...)

	ts := make([]types.Type, 0, len(c.Alternatives())+1)

	for _, a := range c.Alternatives() {
		i, ees, err := i.addVariablesFromPattern(a.Pattern())

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		// Alternative patterns

		tt, ees, err := i.inferType(a.Pattern())

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		ees, err = t.Unify(tt)

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		// Alternative expressions

		tt, ees, err = i.inferType(a.Expression())

		if err != nil {
			return nil, nil, err
		}

		ts = append(ts, tt)
		es = append(es, ees...)
	}

	if d, ok := c.DefaultAlternative(); ok {
		t, ees, err := i.addVariables(map[string]types.Type{d.Variable(): c.Type()}).inferType(
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
}

func (i inferrer) inferLambda(l ast.Lambda) (types.Type, []types.Equation, error) {
	as := make(map[string]types.Type, len(l.Arguments()))

	for _, s := range l.Arguments() {
		as[s] = i.createTypeVariable()
	}

	t, es, err := i.addVariables(as).inferType(l.Expression())

	if err != nil {
		return nil, nil, err
	}

	for i := len(as) - 1; i >= 0; i-- {
		t = types.NewFunction(as[l.Arguments()[i]], t, nil)
	}

	return t, es, nil
}

func (i inferrer) inferLet(l ast.Let) (types.Type, []types.Equation, error) {
	es := []types.Equation{}
	i = i.addVariablesFromBinds(l.Binds())

	for _, b := range l.Binds() {
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

	t, ees, err := i.inferType(l.Expression())

	if err != nil {
		return nil, nil, err
	}

	return t, append(es, ees...), nil
}

func (i inferrer) inferList(l ast.List) (types.Type, []types.Equation, error) {
	t := types.NewList(i.createTypeVariable(), l.Type().DebugInformation())
	es, err := l.Type().Unify(t)

	if err != nil {
		return nil, nil, err
	}

	for _, a := range l.Arguments() {
		tt, ees, err := i.inferType(a.Expression())

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)

		ttt := t.Element()

		if a.Expanded() {
			ttt = t
		}

		ees, err = ttt.Unify(tt)

		if err != nil {
			return nil, nil, err
		}

		es = append(es, ees...)
	}

	return t, es, nil
}

func (i inferrer) inferUnboxed(u ast.Unboxed) (types.Type, []types.Equation, error) {
	t, es, err := i.inferType(u.Content())

	if err != nil {
		return nil, nil, err
	}

	return types.NewUnboxed(t, nil), es, nil
}

func (i inferrer) inferVariable(v ast.Variable) (types.Type, []types.Equation, error) {
	t, ok := i.variables[v.Name()]

	if !ok {
		return nil, nil, fmt.Errorf("variable '%s' not found", v.Name())
	}

	return t, nil, nil
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
			as := make([]ast.Alternative, 0, len(e.Alternatives()))

			for _, a := range e.Alternatives() {
				as = append(
					as,
					ast.NewAlternative(
						a.Pattern().ConvertExpressions(func(e ast.Expression) ast.Expression {
							switch e := e.(type) {
							case ast.List:
								return ast.NewList(i.substituteVariable(e.Type(), ss), e.Arguments())
							}

							return e
						}).(ast.Expression),
						a.Expression(),
					),
				)
			}

			a, ok := e.DefaultAlternative()

			if !ok {
				return ast.NewCaseWithoutDefault(e.Expression(), i.substituteVariable(e.Type(), ss),
					as)
			}

			return ast.NewCase(e.Expression(), i.substituteVariable(e.Type(), ss), as, a)
		case ast.Let:
			bs := make([]ast.Bind, 0, len(e.Binds()))

			for _, b := range e.Binds() {
				bs = append(
					bs,
					ast.NewBind(b.Name(), i.substituteVariable(b.Type(), ss), b.Expression()),
				)
			}

			return ast.NewLet(bs, e.Expression())
		case ast.List:
			return ast.NewList(i.substituteVariable(e.Type(), ss), e.Arguments())
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

func (i inferrer) addVariablesFromPattern(e ast.Expression) (inferrer, []types.Equation, error) {
	l, ok := e.(ast.List)

	if !ok {
		return i, nil, nil
	}

	m := make(map[string]types.Type, len(l.Arguments()))
	es := []types.Equation{}

	for _, a := range l.Arguments() {
		if v, ok := a.Expression().(ast.Variable); ok && a.Expanded() {
			m[v.Name()] = l.Type()
		} else if ok && !a.Expanded() {
			t := types.NewList(i.createTypeVariable(), l.Type().DebugInformation())
			ees, err := l.Type().Unify(t)

			if err != nil {
				return inferrer{}, nil, err
			}

			m[v.Name()] = t.Element()
			es = append(es, ees...)
		} else {
			ii, ees, err := i.addVariablesFromPattern(a.Expression())

			if err != nil {
				return inferrer{}, nil, err
			}

			i = ii
			es = append(es, ees...)
		}
	}

	return i.addVariables(m), es, nil
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
			as := make([]ast.Alternative, 0, len(e.Alternatives()))

			for _, a := range e.Alternatives() {
				as = append(
					as,
					ast.NewAlternative(
						a.Pattern().ConvertExpressions(func(e ast.Expression) ast.Expression {
							switch e := e.(type) {
							case ast.List:
								return ast.NewList(i.createTypeVariable(), e.Arguments())
							}

							return e
						}).(ast.Expression),
						a.Expression(),
					),
				)
			}

			a, ok := e.DefaultAlternative()

			if !ok {
				return ast.NewCaseWithoutDefault(e.Expression(), i.createTypeVariable(), as)
			}

			return ast.NewCase(e.Expression(), i.createTypeVariable(), as, a)
		case ast.Let:
			bs := make([]ast.Bind, 0, len(e.Binds()))

			for _, b := range e.Binds() {
				bs = append(bs, ast.NewBind(b.Name(), i.createTypeVariable(), b.Expression()))
			}

			return ast.NewLet(bs, e.Expression())
		case ast.List:
			return ast.NewList(i.createTypeVariable(), e.Arguments())
		}

		return e
	}).(ast.Module)
}
