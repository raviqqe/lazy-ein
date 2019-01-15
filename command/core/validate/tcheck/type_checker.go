package tcheck

import (
	"errors"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/types"
)

type typeChecker struct {
	variables map[string]types.Type
}

func newTypeChecker() typeChecker {
	return typeChecker{map[string]types.Type{}}
}

func (c typeChecker) Check(m ast.Module) error {
	_, err := c.checkBinds(m.Binds())
	return err
}

func (c typeChecker) checkBinds(bs []ast.Bind) (typeChecker, error) {
	c = c.addBinds(bs)

	for _, b := range bs {
		if err := c.checkLambda(b.Lambda()); err != nil {
			return typeChecker{}, err
		}
	}

	return c, nil
}

func (c typeChecker) checkLambda(l ast.Lambda) error {
	t, err := c.addArguments(l.Arguments()).checkExpression(l.Body())

	if err != nil {
		return err
	} else if !types.Equal(l.ResultType(), t) {
		return errors.New("unmatched result types")
	}

	return nil
}

func (c typeChecker) checkExpression(e ast.Expression) (types.Type, error) {
	switch e := e.(type) {
	case ast.AlgebraicCase:
		t, err := c.checkExpression(e.Argument())

		if err != nil {
			return nil, err
		} else if len(e.Alternatives()) == 0 {
			d, _ := e.DefaultAlternative()
			return c.checkDefaultAlternative(d, t)
		}

		tt, err := c.checkAlgebraicAlternative(e.Alternatives()[0])

		if err != nil {
			return nil, err
		}

		for _, a := range e.Alternatives()[1:] {
			if err := c.checkTypes(a.Constructor().AlgebraicType(), types.Unbox(t)); err != nil {
				return nil, err
			}

			if ttt, err := c.checkAlgebraicAlternative(a); err != nil {
				return nil, err
			} else if err := c.checkTypes(ttt, tt); err != nil {
				return nil, err
			}
		}

		d, ok := e.DefaultAlternative()

		if !ok {
			return tt, nil
		}

		return c.checkDefaultAlternative(d, t)
	case ast.Atom:
		return c.inferAtomType(e)
	case ast.ConstructorApplication:
		if len(e.Arguments()) != len(e.Constructor().ConstructorType().Elements()) {
			return nil, errors.New("invalid number of constructor arguments")
		}

		for i, a := range e.Arguments() {
			t, err := c.inferAtomType(a)

			if err != nil {
				return nil, err
			}

			tt := e.Constructor().ConstructorType().Elements()[i]

			if err := c.checkTypes(t, tt); err != nil {
				return nil, err
			}
		}

		return e.Constructor().AlgebraicType(), nil
	case ast.FunctionApplication:
		if len(e.Arguments()) == 0 {
			return c.getVariableType(e.Function())
		}

		t, err := c.getVariableType(e.Function())

		if err != nil {
			return nil, err
		}

		f, ok := t.(types.Function)

		if !ok {
			return nil, errors.New("not a function")
		} else if len(f.Arguments()) != len(e.Arguments()) {
			return nil, errors.New("invalid number of function arguments")
		}

		for i, t := range f.Arguments() {
			tt, err := c.inferAtomType(e.Arguments()[i])

			if err != nil {
				return nil, err
			} else if err := c.checkTypes(t, tt); err != nil {
				return nil, err
			}
		}

		return f.Result(), nil
	case ast.Let:
		c, err := c.checkBinds(e.Binds())

		if err != nil {
			return nil, err
		}

		return c.checkExpression(e.Expression())
	case ast.PrimitiveCase:
		if t, err := c.checkExpression(e.Argument()); err != nil {
			return nil, err
		} else if err := c.checkTypes(t, e.Type()); err != nil {
			return nil, err
		} else if d, ok := e.DefaultAlternative(); ok && len(e.Alternatives()) == 0 {
			return c.checkDefaultAlternative(d, e.Type())
		}

		t, err := c.checkExpression(e.Alternatives()[0].Expression())

		if err != nil {
			return nil, err
		}

		for _, a := range e.Alternatives()[1:] {
			if err := c.checkTypes(c.getLiteralType(a.Literal()), types.Unbox(e.Type())); err != nil {
				return nil, err
			}

			if tt, err := c.checkExpression(a.Expression()); err != nil {
				return nil, err
			} else if err := c.checkTypes(t, tt); err != nil {
				return nil, err
			}
		}

		d, ok := e.DefaultAlternative()

		if !ok {
			return t, nil
		}

		return c.checkDefaultAlternative(d, e.Type())
	case ast.PrimitiveOperation:
		for _, a := range e.Arguments() {
			t, err := c.inferAtomType(a)

			if err != nil {
				return nil, err
			} else if err := c.checkTypes(t, types.NewFloat64()); err != nil {
				return nil, err
			}
		}

		return types.NewFloat64(), nil
	}

	panic("unreachable")
}

func (c typeChecker) checkAlgebraicAlternative(a ast.AlgebraicAlternative) (types.Type, error) {
	return c.addConstructorElements(a.Constructor(), a.ElementNames()).checkExpression(
		a.Expression(),
	)
}

func (c typeChecker) checkDefaultAlternative(d ast.DefaultAlternative, t types.Type) (types.Type, error) {
	return c.addVariable(d.Variable(), types.Unbox(t)).checkExpression(d.Expression())
}

func (typeChecker) checkTypes(t, tt types.Type) error {
	if !types.Equal(t, tt) {
		return errors.New("unmatched types")
	}

	return nil
}

func (c typeChecker) inferAtomType(a ast.Atom) (types.Type, error) {
	switch a := a.(type) {
	case ast.Literal:
		return c.getLiteralType(a), nil
	case ast.Variable:
		return c.getVariableType(a)
	}

	panic("unreachable")
}

func (c typeChecker) getLiteralType(l ast.Literal) types.Type {
	switch l.(type) {
	case ast.Literal:
		return types.NewFloat64()
	}

	panic("unreachable")
}

func (c typeChecker) getVariableType(v ast.Variable) (types.Type, error) {
	t, ok := c.variables[v.Name()]

	if !ok {
		return nil, errors.New("variable not found")
	}

	return t, nil
}

func (c typeChecker) addArguments(as []ast.Argument) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+len(as))

	for k, v := range c.variables {
		vs[k] = v
	}

	for _, a := range as {
		vs[a.Name()] = a.Type()
	}

	return typeChecker{vs}
}

func (c typeChecker) addBinds(bs []ast.Bind) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+len(bs))

	for k, v := range c.variables {
		vs[k] = v
	}

	for _, b := range bs {
		vs[b.Name()] = b.Lambda().Type()
	}

	return typeChecker{vs}
}

func (c typeChecker) addConstructorElements(cc ast.Constructor, ss []string) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+len(ss))

	for k, v := range c.variables {
		vs[k] = v
	}

	for i, s := range ss {
		if s != "" {
			vs[s] = cc.ConstructorType().Elements()[i]
		}
	}

	return typeChecker{vs}
}

func (c typeChecker) addVariable(s string, t types.Type) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+1)

	for k, v := range c.variables {
		vs[k] = v
	}

	vs[s] = t

	return typeChecker{vs}
}
