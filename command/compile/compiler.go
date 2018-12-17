package compile

import (
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/types"
)

type compiler struct {
	variables       map[string]coretypes.Type
	globalVariables map[string]struct{}
}

func newCompiler(m ast.Module) (compiler, error) {
	vs := make(map[string]coretypes.Type, len(m.Binds()))
	gs := make(map[string]struct{}, len(m.Binds()))

	for _, b := range m.Binds() {
		t, err := types.Box(b.Type()).ToCore()

		if err != nil {
			return compiler{}, err
		}

		vs[b.Name()] = t
		gs[b.Name()] = struct{}{}
	}

	return compiler{vs, gs}, nil
}

func (c compiler) Compile(m ast.Module) (coreast.Module, error) {
	bs := make([]coreast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		b, err := c.compileBind(b, true)

		if err != nil {
			return coreast.Module{}, err
		}

		bs = append(bs, b)
	}

	return coreast.NewModule(m.Name(), nil, bs), nil
}

func (c compiler) compileBind(b ast.Bind, global bool) (coreast.Bind, error) {
	t, err := b.Type().ToCore()

	if err != nil {
		return coreast.Bind{}, err
	}

	c = c.addVariable(b.Name(), t)

	if len(b.Arguments()) == 0 {
		vs, err := c.compileFreeVariables(b, global)

		if err != nil {
			return coreast.Bind{}, err
		}

		e, err := c.compileExpression(b.Expression())

		if err != nil {
			return coreast.Bind{}, err
		}

		return coreast.NewBind(b.Name(), coreast.NewLambda(vs, true, nil, e, t)), nil
	}

	f, ok := t.(coretypes.Function)

	if !ok {
		return coreast.Bind{}, types.NewTypeError("not a function", b.Type().DebugInformation())
	}

	as := make([]coreast.Argument, 0, len(b.Arguments()))

	for i, n := range b.Arguments() {
		t := f.Arguments()[i]
		as = append(as, coreast.NewArgument(n, t))
		c = c.addVariable(n, t)
	}

	vs, err := c.compileFreeVariables(b, global)

	if err != nil {
		return coreast.Bind{}, err
	}

	e, err := c.compileExpression(b.Expression())

	if err != nil {
		return coreast.Bind{}, err
	}

	return coreast.NewBind(b.Name(), coreast.NewLambda(vs, false, as, e, f.Result())), nil
}

func (c compiler) compileExpression(e ast.Expression) (coreast.Expression, error) {
	switch e := e.(type) {
	case ast.Application:
		as := make([]coreast.Atom, 0, len(e.Arguments()))

		for _, a := range e.Arguments() {
			as = append(as, coreast.NewVariable(a.(ast.Variable).Name()))
		}

		return coreast.NewApplication(
			coreast.NewVariable(e.Function().(ast.Variable).Name()),
			as,
		), nil
	case ast.Let:
		bs := make([]coreast.Bind, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			b, err := c.compileBind(b, false)

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
	case ast.Unboxed:
		// TODO: Handle other literals.
		switch e := e.Content().(type) {
		case ast.Number:
			return coreast.NewFloat64(e.Value()), nil
		}
	case ast.Variable:
		return coreast.NewApplication(coreast.NewVariable(e.Name()), nil), nil
	}

	panic("unreahable")
}

func (c compiler) compileFreeVariables(b ast.Bind, global bool) ([]coreast.Argument, error) {
	if global {
		return nil, nil
	}

	ss := c.findFreeVariables(b)

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

func (c compiler) findFreeVariables(b ast.Bind) []string {
	m := make(map[string]struct{}, len(c.globalVariables)+len(b.Arguments()))

	for k := range c.globalVariables {
		m[k] = struct{}{}
	}

	for _, k := range b.Arguments() {
		m[k] = struct{}{}
	}

	return newFreeVariableFinder(m).Find(b.Expression())
}

func (c compiler) addVariable(s string, t coretypes.Type) compiler {
	m := make(map[string]coretypes.Type, len(c.variables)+1)

	for k, v := range c.variables {
		m[k] = v
	}

	m[s] = t

	return compiler{m, c.globalVariables}
}
