package tinfer

import (
	"errors"
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
)

type inferrer struct {
	variables map[string]types.Type
}

func newInferrer(m ast.Module) inferrer {
	vs := make(map[string]types.Type, len(m.Binds()))

	for _, b := range m.Binds() {
		vs[b.Name()] = types.Box(b.Type())
	}

	return inferrer{vs}
}

func (i inferrer) Infer(m ast.Module) (ast.Module, error) {
	bs := make([]ast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		vs := make(map[string]types.Type)
		t := b.Type()

		for _, s := range b.Arguments() {
			f, ok := t.(types.Function)

			if !ok {
				return ast.Module{}, errors.New("too many arguments")
			}

			vs[s] = f.Argument()

			t = f.Result()
		}

		i = i.addVariables(vs)

		e, err := i.inferExpression(b.Expression())

		if err != nil {
			return ast.Module{}, err
		}

		bs = append(bs, ast.NewBind(b.Name(), b.Arguments(), b.Type(), e))
	}

	return ast.NewModule(m.Name(), bs), nil
}

func (i inferrer) inferExpression(e ast.Expression) (ast.Expression, error) {
	l, ok := e.(ast.Let)

	if !ok {
		return e, nil
	}

	i = i.addVariablesFromBinds(l.Binds())

	bs := make([]ast.Bind, 0, len(l.Binds()))

	for _, b := range l.Binds() {
		t := b.Type()

		if len(b.Arguments()) > 0 {
			vs := make(map[string]types.Type)
			t = types.Type(types.NewVariable(b.Type().DebugInformation()))
			f := t

			for i := len(b.Arguments()) - 1; i >= 0; i-- {
				v := types.NewVariable(nil)
				vs[b.Arguments()[i]] = v
				f = types.NewFunction(v, f, b.Type().DebugInformation())
			}

			i = i.addVariables(vs)

			if err := b.Type().Unify(f); err != nil {
				return nil, err
			}
		}

		e, err := i.inferExpression(b.Expression())

		if err != nil {
			return nil, err
		}

		tt, err := i.inferExpressionType(e)

		if err != nil {
			return nil, err
		} else if err = t.Unify(tt); err != nil {
			return nil, err
		}

		bs = append(bs, ast.NewBind(b.Name(), b.Arguments(), b.Type(), e))
	}

	e, err := i.inferExpression(l.Expression())

	if err != nil {
		return nil, err
	}

	return ast.NewLet(bs, e), nil
}

func (i inferrer) inferExpressionType(e ast.Expression) (types.Type, error) {
	switch e := e.(type) {
	case ast.Application:
		t, err := i.inferExpressionType(e.Function())

		if err != nil {
			return nil, err
		}

		for _, a := range e.Arguments() {
			f, ok := t.(types.Function)

			if !ok {
				return nil, types.NewTypeError("not a function", t.DebugInformation())
			}

			a, err := i.inferExpressionType(a)

			if err != nil {
				return nil, err
			} else if err := f.Argument().Unify(a); err != nil {
				return nil, err
			}

			t = f.Result()
		}

		return t, nil
	case ast.Let:
		vs := make(map[string]types.Type, len(e.Binds()))

		for _, b := range e.Binds() {
			vs[b.Name()] = b.Type()
		}

		return i.addVariables(vs).inferExpressionType(e.Expression())
	case ast.Number:
		return types.NewNumber(nil), nil
	case ast.Variable:
		t, ok := i.variables[e.Name()]

		if !ok {
			return nil, fmt.Errorf("variable '%s' not found", e.Name())
		}

		return t, nil
	}

	panic("unreachable")
}

func (i inferrer) addVariablesFromBinds(bs []ast.Bind) inferrer {
	m := make(map[string]types.Type, len(i.variables)+len(bs))

	for k, v := range i.variables {
		m[k] = v
	}

	for _, b := range bs {
		m[b.Name()] = b.Type()
	}

	return inferrer{m}
}

func (i inferrer) addVariables(vs map[string]types.Type) inferrer {
	m := make(map[string]types.Type, len(i.variables)+len(vs))

	for k, v := range i.variables {
		m[k] = v
	}

	for k, v := range vs {
		m[k] = v
	}

	return inferrer{m}
}
