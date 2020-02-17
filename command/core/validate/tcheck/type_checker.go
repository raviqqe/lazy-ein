package tcheck

import (
	"errors"

	"github.com/raviqqe/lazy-ein/command/core/ast"
	"github.com/raviqqe/lazy-ein/command/core/types"
)

type typeChecker struct {
	variables map[string]types.Type
}

func newTypeChecker() typeChecker {
	return typeChecker{map[string]types.Type{}}
}

func (c typeChecker) Check(m ast.Module) error {
	_, err := c.addDeclarations(m.Declarations()).checkBinds(m.Binds())
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
	t, err := c.addArguments(l).checkExpression(l.Body())

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
		return c.checkAlgebraicCase(e)
	case ast.Atom:
		return c.checkAtom(e)
	case ast.ConstructorApplication:
		return c.checkConstructorApplication(e)
	case ast.FunctionApplication:
		return c.checkFunctionApplication(e)
	case ast.Let:
		return c.checkLet(e)
	case ast.PrimitiveCase:
		return c.checkPrimitiveCase(e)
	case ast.PrimitiveOperation:
		return c.checkPrimitiveOperation(e)
	}

	panic("unreachable")
}

func (c typeChecker) checkAlgebraicCase(ac ast.AlgebraicCase) (types.Type, error) {
	t, err := c.checkExpression(ac.Argument())

	if err != nil {
		return nil, err
	} else if len(ac.Alternatives()) == 0 {
		d, _ := ac.DefaultAlternative()
		return c.checkDefaultAlternative(d, t)
	}

	tt, err := c.checkAlgebraicAlternative(ac.Alternatives()[0])

	if err != nil {
		return nil, err
	}

	for _, a := range ac.Alternatives()[1:] {
		if err := c.checkTypes(a.Constructor().AlgebraicType(), types.Unbox(t)); err != nil {
			return nil, err
		}

		if ttt, err := c.checkAlgebraicAlternative(a); err != nil {
			return nil, err
		} else if err := c.checkTypes(ttt, tt); err != nil {
			return nil, err
		}
	}

	d, ok := ac.DefaultAlternative()

	if !ok {
		return tt, nil
	}

	return c.checkDefaultAlternative(d, t)
}

func (c typeChecker) checkAlgebraicAlternative(a ast.AlgebraicAlternative) (types.Type, error) {
	return c.addConstructorElements(a.Constructor(), a.ElementNames()).checkExpression(
		a.Expression(),
	)
}

func (c typeChecker) checkConstructorApplication(a ast.ConstructorApplication) (types.Type, error) {
	if len(a.Arguments()) != len(a.Constructor().ConstructorType().Elements()) {
		return nil, errors.New("invalid number of constructor arguments")
	}

	for i, arg := range a.Arguments() {
		t, err := c.checkAtom(arg)

		if err != nil {
			return nil, err
		}

		tt := a.Constructor().ConstructorType().Elements()[i]

		if err := c.checkTypes(t, tt); err != nil {
			return nil, err
		}
	}

	return a.Constructor().AlgebraicType(), nil
}

func (c typeChecker) checkFunctionApplication(a ast.FunctionApplication) (types.Type, error) {
	if len(a.Arguments()) == 0 {
		return c.checkVariable(a.Function())
	}

	t, err := c.checkVariable(a.Function())

	if err != nil {
		return nil, err
	}

	f, ok := t.(types.Function)

	if !ok {
		return nil, errors.New("not a function")
	} else if len(f.Arguments()) != len(a.Arguments()) {
		return nil, errors.New("invalid number of function arguments")
	}

	for i, t := range f.Arguments() {
		tt, err := c.checkAtom(a.Arguments()[i])

		if err != nil {
			return nil, err
		} else if err := c.checkTypes(t, tt); err != nil {
			return nil, err
		}
	}

	return f.Result(), nil
}

func (c typeChecker) checkLet(l ast.Let) (types.Type, error) {
	c, err := c.checkBinds(l.Binds())

	if err != nil {
		return nil, err
	}

	return c.checkExpression(l.Expression())
}

func (c typeChecker) checkPrimitiveCase(pc ast.PrimitiveCase) (types.Type, error) {
	if t, err := c.checkExpression(pc.Argument()); err != nil {
		return nil, err
	} else if err := c.checkTypes(t, pc.Type()); err != nil {
		return nil, err
	} else if d, ok := pc.DefaultAlternative(); ok && len(pc.Alternatives()) == 0 {
		return c.checkDefaultAlternative(d, pc.Type())
	}

	t, err := c.checkExpression(pc.Alternatives()[0].Expression())

	if err != nil {
		return nil, err
	}

	for _, a := range pc.Alternatives()[1:] {
		if err := c.checkTypes(c.getLiteralType(a.Literal()), types.Unbox(pc.Type())); err != nil {
			return nil, err
		}

		if tt, err := c.checkExpression(a.Expression()); err != nil {
			return nil, err
		} else if err := c.checkTypes(t, tt); err != nil {
			return nil, err
		}
	}

	d, ok := pc.DefaultAlternative()

	if !ok {
		return t, nil
	}

	return c.checkDefaultAlternative(d, pc.Type())
}

func (c typeChecker) checkPrimitiveOperation(o ast.PrimitiveOperation) (types.Type, error) {
	for _, a := range o.Arguments() {
		t, err := c.checkAtom(a)

		if err != nil {
			return nil, err
		} else if err := c.checkTypes(t, types.NewFloat64()); err != nil {
			return nil, err
		}
	}

	return types.NewFloat64(), nil
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

func (c typeChecker) checkAtom(a ast.Atom) (types.Type, error) {
	switch a := a.(type) {
	case ast.Literal:
		return c.getLiteralType(a), nil
	case ast.Variable:
		return c.checkVariable(a)
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

func (c typeChecker) checkVariable(v ast.Variable) (types.Type, error) {
	t, ok := c.variables[v.Name()]

	if !ok {
		return nil, errors.New("variable not found")
	}

	return t, nil
}

func (c typeChecker) addArguments(l ast.Lambda) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+len(l.ArgumentNames()))

	for k, v := range c.variables {
		vs[k] = v
	}

	for i, s := range l.ArgumentNames() {
		vs[s] = l.ArgumentTypes()[i]
	}

	return typeChecker{vs}
}

func (c typeChecker) addDeclarations(ds []ast.Declaration) typeChecker {
	vs := make(map[string]types.Type, len(c.variables)+len(ds))

	for k, v := range c.variables {
		vs[k] = v
	}

	for _, d := range ds {
		vs[d.Name()] = d.Lambda().Type()
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
