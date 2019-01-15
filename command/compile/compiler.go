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
		vs[b.Name()] = types.Box(b.Type()).ToCore()
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

		bs = append(bs, coreast.NewBind(b.Name(), b.Lambda().ClearFreeVariables()))
	}

	return coreast.NewModule(m.Name(), bs), nil
}

func (c compiler) compileBind(b ast.Bind) (coreast.Bind, error) {
	t := b.Type().ToCore()
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

		return coreast.NewBind(
			b.Name(),
			coreast.NewVariableLambda(vs, true, e, t.(coretypes.Boxable)),
		), nil
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

	return coreast.NewBind(b.Name(), coreast.NewFunctionLambda(vs, as, e, f.Result())), nil
}

func (c compiler) compileExpression(e ast.Expression) (coreast.Expression, error) {
	switch e := e.(type) {
	case ast.Application:
		as := make([]coreast.Atom, 0, len(e.Arguments()))

		for _, a := range e.Arguments() {
			as = append(as, coreast.NewVariable(a.(ast.Variable).Name()))
		}

		return coreast.NewFunctionApplication(
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

		x := l.(coreast.FunctionApplication).Function().Name()
		y := r.(coreast.FunctionApplication).Function().Name()

		vs, err := c.compileFreeVariables(e)

		if err != nil {
			return nil, err
		}

		return coreast.NewLet(
			[]coreast.Bind{
				coreast.NewBind(
					"$boxedResult",
					coreast.NewVariableLambda(
						vs,
						true,
						c.bindNumberPrimitive(
							coreast.NewFunctionApplication(coreast.NewVariable(x), nil),
							"$lhs",
							c.bindNumberPrimitive(
								coreast.NewFunctionApplication(coreast.NewVariable(y), nil),
								"$rhs",
								coreast.NewPrimitiveCase(
									coreast.NewPrimitiveOperation(
										binaryOperatorToPrimitive(e.Operator()),
										[]coreast.Atom{
											coreast.NewVariable("$lhs"),
											coreast.NewVariable("$rhs"),
										},
									),
									coretypes.NewFloat64(),
									nil,
									coreast.NewDefaultAlternative(
										"$result",
										coreast.NewConstructorApplication(
											types.NewNumber(nil).CoreConstructor(),
											[]coreast.Atom{coreast.NewVariable("$result")},
										),
									),
								),
							),
						),
						coretypes.Unbox(types.NewNumber(nil).ToCore()).(coretypes.Algebraic),
					),
				),
			},
			coreast.NewFunctionApplication(coreast.NewVariable("$boxedResult"), nil),
		), nil
	case ast.Case:
		ee, err := c.compileExpression(e.Expression())

		if err != nil {
			return nil, err
		}

		as := make([]coreast.PrimitiveAlternative, 0, len(e.Alternatives()))

		for _, a := range e.Alternatives() {
			e, err := c.compileExpression(a.Expression())

			if err != nil {
				return nil, err
			}

			as = append(as, coreast.NewPrimitiveAlternative(c.compileUnboxedLiteral(a.Literal()), e))
		}

		d, ok := e.DefaultAlternative()

		if !ok {
			return coreast.NewPrimitiveCaseWithoutDefault(
				c.extractNumberPrimitive(ee),
				coretypes.NewFloat64(),
				as,
			), nil
		}

		c = c.addVariable(d.Variable(), e.Type().ToCore())
		de, err := c.compileExpression(d.Expression())

		if err != nil {
			return nil, err
		}

		vs, err := c.compileFreeVariables(e.Expression())

		if err != nil {
			return nil, err
		}

		return coreast.NewLet(
			[]coreast.Bind{
				coreast.NewBind(
					d.Variable(),
					coreast.NewVariableLambda(vs, true, ee, e.Type().ToCore().(coretypes.Boxable)),
				),
			},
			coreast.NewPrimitiveCase(
				c.extractNumberPrimitive(
					coreast.NewFunctionApplication(coreast.NewVariable(d.Variable()), nil),
				),
				coretypes.NewFloat64(),
				as,
				coreast.NewDefaultAlternative("", de),
			),
		), nil
	case ast.Let:
		bs := make([]coreast.Bind, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			c = c.addVariable(b.Name(), b.Type().ToCore())
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
			return coreast.NewConstructorApplication(
				types.NewNumber(nil).CoreConstructor(),
				[]coreast.Atom{coreast.NewFloat64(e.Value())},
			), nil
		}
	case ast.Variable:
		return coreast.NewFunctionApplication(coreast.NewVariable(e.Name()), nil), nil
	}

	panic("unreahable")
}

func (compiler) compileUnboxedLiteral(l ast.Literal) coreast.Literal {
	switch l := l.(type) {
	case ast.Number:
		return coreast.NewFloat64(l.Value())
	}

	panic("unreahable")
}

func (c compiler) extractNumberPrimitive(e coreast.Expression) coreast.Expression {
	return c.bindNumberPrimitive(
		e,
		"$primitive",
		coreast.NewFunctionApplication(coreast.NewVariable("$primitive"), nil),
	)
}

func (compiler) bindNumberPrimitive(
	e coreast.Expression,
	s string,
	ee coreast.Expression,
) coreast.Expression {
	return coreast.NewAlgebraicCaseWithoutDefault(
		e,
		[]coreast.AlgebraicAlternative{
			coreast.NewAlgebraicAlternative(
				types.NewNumber(nil).CoreConstructor(),
				[]string{s},
				ee,
			),
		},
	)
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
