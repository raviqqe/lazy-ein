package compile

import (
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	coreast "github.com/ein-lang/ein/command/core/ast"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"github.com/ein-lang/ein/command/types"
)

type compiler struct {
	variables          map[string]coretypes.Type
	freeVariableFinder freeVariableFinder
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

	return compiler{vs, newFreeVariableFinder(gs)}, nil
}

func (c compiler) Compile(m ast.Module) (coreast.Module, error) {
	bs := make([]coreast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		b, err := c.compileBind(b)

		if err != nil {
			return coreast.Module{}, err
		}

		l := b.Lambda()

		bs = append(bs,
			coreast.NewBind(
				b.Name(),
				coreast.NewLambda(nil, l.IsUpdatable(), l.Arguments(), l.Body(), l.ResultType()),
			),
		)
	}

	return coreast.NewModule(m.Name(), nil, bs), nil
}

func (c compiler) compileBind(b ast.Bind) (coreast.Bind, error) {
	t, err := b.Type().ToCore()

	if err != nil {
		return coreast.Bind{}, err
	}

	c = c.addVariable(b.Name(), t)
	l, ok := b.Expression().(ast.Lambda)

	if !ok {
		vs, err := c.compileFreeVariables(b.Expression())

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

	as := make([]coreast.Argument, 0, len(l.Arguments()))

	for i, s := range l.Arguments() {
		t := f.Arguments()[i]
		as = append(as, coreast.NewArgument(s, t))
		c = c.addVariable(s, t)
	}

	vs, err := c.compileFreeVariables(b.Expression())

	if err != nil {
		return coreast.Bind{}, err
	}

	e, err := c.compileExpression(l.Expression())

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
	case ast.BinaryOperation:
		l, err := c.compileExpression(e.LHS())

		if err != nil {
			return nil, err
		}

		r, err := c.compileExpression(e.RHS())

		if err != nil {
			return nil, err
		}

		x := l.(coreast.Application).Function().Name()
		y := r.(coreast.Application).Function().Name()

		vs, err := c.compileFreeVariables(e)

		if err != nil {
			return nil, err
		}

		return coreast.NewLet(
			[]coreast.Bind{
				coreast.NewBind(
					"result",
					coreast.NewLambda(
						vs,
						true,
						nil,
						coreast.NewPrimitiveCase(
							coreast.NewApplication(coreast.NewVariable(x), nil),
							coretypes.NewBoxed(coretypes.NewFloat64()),
							nil,
							coreast.NewDefaultAlternative(
								"lhs",
								coreast.NewPrimitiveCase(
									coreast.NewApplication(coreast.NewVariable(y), nil),
									coretypes.NewBoxed(coretypes.NewFloat64()),
									nil,
									coreast.NewDefaultAlternative(
										"rhs",
										coreast.NewPrimitiveOperation(
											binaryOperatorToPrimitive(e.Operator()),
											[]coreast.Atom{
												coreast.NewVariable("lhs"),
												coreast.NewVariable("rhs"),
											},
										),
									),
								),
							),
						),
						coretypes.NewFloat64(),
					),
				),
			},
			coreast.NewApplication(coreast.NewVariable("result"), nil),
		), nil
	case ast.Let:
		bs := make([]coreast.Bind, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			t, err := b.Type().ToCore()

			if err != nil {
				return nil, err
			}

			c = c.addVariable(b.Name(), t)
		}

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
