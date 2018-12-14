package compile

import (
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
)

type compiler struct {
	variables          map[string]coretypes.Type
	freeVariableFinder freeVariableFinder
}

func newCompiler(m ast.Module) compiler {
	vs := make(map[string]coretypes.Type, len(m.Binds()))
	gs := make(map[string]struct{}, len(m.Binds()))

	for _, b := range m.Binds() {
		vs[b.Name()] = b.Type().ToCore()
		gs[b.Name()] = struct{}{}
	}

	return compiler{vs, newFreeVariableFinder(gs)}
}

func (c compiler) Compile(m ast.Module) (coreast.Module, error) {
	bs := make([]coreast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		b, err := c.compileBind(b)

		if err != nil {
			return coreast.Module{}, err
		}

		bs = append(bs, b)
	}

	return coreast.NewModule(m.Name(), nil, bs), nil
}

func (c compiler) compileBind(b ast.Bind) (coreast.Bind, error) {
	c = c.addVariable(b.Name(), b.Type().ToCore())

	if len(b.Arguments()) == 0 {
		vs, err := c.compileFreeVariables(b.Expression())

		if err != nil {
			return coreast.Bind{}, err
		}

		e, err := c.compileExpression(b.Expression())

		if err != nil {
			return coreast.Bind{}, err
		}

		return coreast.NewBind(
			b.Name(),
			coreast.NewLambda(vs, true, nil, e, b.Type().ToCore()),
		), nil
	}

	t := b.Type().ToCore().(coretypes.Function)
	as := make([]coreast.Argument, 0, len(b.Arguments()))

	for i, n := range b.Arguments() {
		t := t.Arguments()[i]
		as = append(as, coreast.NewArgument(n, t))
		c = c.addVariable(n, t)
	}

	vs, err := c.compileFreeVariables(b.Expression())

	if err != nil {
		return coreast.Bind{}, err
	}

	e, err := c.compileExpression(b.Expression())

	if err != nil {
		return coreast.Bind{}, err
	}

	return coreast.NewBind(
		b.Name(),
		coreast.NewLambda(vs, false, as, e, t.Result()),
	), nil
}

func (c compiler) compileExpression(e ast.Expression) (coreast.Expression, error) {
	switch e := e.(type) {
	case ast.Let:
		bs := make([]coreast.Bind, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			b, err := c.compileBind(b)

			if err != nil {
				return nil, err
			}

			bs = append(bs, b)
		}

		ee, err := c.compileExpression(e.Expression())

		if err != nil {
			return nil, err
		}

		return coreast.NewLet(bs, ee), nil
	case ast.Number:
		return coreast.NewFloat64(e.Value()), nil
	case ast.Variable:
		return coreast.NewApplication(coreast.NewVariable(e.Name()), nil), nil
	}

	panic("unreahable")
}

func (c compiler) compileFreeVariables(e ast.Expression) ([]coreast.Argument, error) {
	ss := c.freeVariableFinder.Find(e)

	if len(ss) == 0 {
		return nil, nil // to normalize empty slices
	}

	as := make([]coreast.Argument, 0, len(ss))

	for _, s := range ss {
		t, ok := c.variables[s]

		if !ok {
			return nil, fmt.Errorf("variable '%s' not found", s)
		}

		as = append(as, coreast.NewArgument(s, t))
	}

	return as, nil
}

func (c compiler) addVariable(s string, t coretypes.Type) compiler {
	m := make(map[string]coretypes.Type, len(c.variables)+1)

	for k, v := range c.variables {
		m[k] = v
	}

	m[s] = t

	return compiler{m, c.freeVariableFinder}
}
