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
			coreast.NewVariableLambda(vs, true, e, t.(coretypes.Bindable)),
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
		return c.compileApplication(e)
	case ast.BinaryOperation:
		return c.compileBinaryOperation(e)
	case ast.Case:
		return c.compileCase(e)
	case ast.Let:
		return c.compileLet(e)
	case ast.List:
		return c.compileList(e)
	case ast.Unboxed:
		return c.compileUnboxed(e)
	case ast.Variable:
		return coreast.NewFunctionApplication(coreast.NewVariable(e.Name()), nil), nil
	}

	panic("unreachable")
}

func (c compiler) compileApplication(a ast.Application) (coreast.Expression, error) {
	as := make([]coreast.Atom, 0, len(a.Arguments()))

	for _, a := range a.Arguments() {
		as = append(as, coreast.NewVariable(a.(ast.Variable).Name()))
	}

	return coreast.NewFunctionApplication(
		coreast.NewVariable(a.Function().(ast.Variable).Name()),
		as,
	), nil
}

func (c compiler) compileBinaryOperation(o ast.BinaryOperation) (coreast.Expression, error) {
	l, err := c.compileExpression(o.LHS())

	if err != nil {
		return nil, err
	}

	r, err := c.compileExpression(o.RHS())

	if err != nil {
		return nil, err
	}

	x := l.(coreast.FunctionApplication).Function().Name()
	y := r.(coreast.FunctionApplication).Function().Name()

	vs, err := c.compileFreeVariables(o)

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
									binaryOperatorToPrimitive(o.Operator()),
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
}

func (c compiler) compileCase(cc ast.Case) (coreast.Expression, error) {
	switch cc.Type().(type) {
	case types.Number:
		return c.compilePrimitiveCase(cc)
	case types.List:
		return newListCaseCompiler(c).Compile(cc)
	}

	panic("unreachable")
}

func (c compiler) compilePrimitiveCase(cc ast.Case) (coreast.Expression, error) {
	arg, err := c.compileExpression(cc.Expression())

	if err != nil {
		return nil, err
	}

	as := make([]coreast.PrimitiveAlternative, 0, len(cc.Alternatives()))

	for _, a := range cc.Alternatives() {
		e, err := c.compileExpression(a.Expression())

		if err != nil {
			return nil, err
		}

		as = append(
			as,
			coreast.NewPrimitiveAlternative(c.compileUnboxedLiteral(a.Pattern().(ast.Literal)), e),
		)
	}

	d, ok := cc.DefaultAlternative()

	if !ok {
		return coreast.NewPrimitiveCaseWithoutDefault(
			c.extractNumberPrimitive(arg),
			coretypes.NewFloat64(),
			as,
		), nil
	}

	c = c.addVariable(d.Variable(), cc.Type().ToCore())
	de, err := c.compileExpression(d.Expression())

	if err != nil {
		return nil, err
	}

	vs, err := c.compileFreeVariables(cc.Expression())

	if err != nil {
		return nil, err
	}

	return coreast.NewLet(
		[]coreast.Bind{
			coreast.NewBind(
				d.Variable(),
				coreast.NewVariableLambda(vs, true, arg, cc.Type().ToCore().(coretypes.Bindable)),
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
}

func (c compiler) compileLet(l ast.Let) (coreast.Expression, error) {
	bs := make([]coreast.Bind, 0, len(l.Binds()))

	for _, b := range l.Binds() {
		c = c.addVariable(b.Name(), b.Type().ToCore())
	}

	for _, b := range l.Binds() {
		b, err := c.compileBind(b)

		if err != nil {
			return nil, err
		}

		bs = append(bs, b)
	}

	ee, err := c.compileExpression(l.Expression())

	if err != nil {
		return nil, err
	}

	return coreast.NewLet(bs, ee), nil
}

func (c compiler) compileList(l ast.List) (coreast.Expression, error) {
	t := coretypes.Unbox(l.Type().ToCore()).(coretypes.Algebraic)
	s := "$nil"

	bs := make([]coreast.Bind, 0, len(l.Arguments())+1)
	bs = append(
		bs,
		coreast.NewBind(
			s,
			coreast.NewVariableLambda(
				nil,
				true,
				coreast.NewConstructorApplication(coreast.NewConstructor(t, 1), nil),
				t,
			),
		),
	)

	for i := range l.Arguments() {
		a := l.Arguments()[len(l.Arguments())-1-i]

		if a.Expanded() {
			panic("not implemented")
		}

		e := a.Expression().(ast.Variable)
		ss := fmt.Sprintf("$list-%v", i)
		vs := []coreast.Argument{coreast.NewArgument(s, coretypes.NewBoxed(t))}

		if len(c.freeVariableFinder.Find(e)) > 0 {
			vs = append(vs, coreast.NewArgument(e.Name(), t.Constructors()[0].Elements()[i]))
		}

		bs = append(
			bs,
			coreast.NewBind(
				ss,
				coreast.NewVariableLambda(
					vs,
					true,
					coreast.NewConstructorApplication(
						coreast.NewConstructor(t, 0),
						[]coreast.Atom{
							coreast.NewVariable(e.Name()),
							coreast.NewVariable(s),
						},
					),
					t,
				),
			),
		)

		s = ss
	}

	return coreast.NewLet(bs, coreast.NewFunctionApplication(coreast.NewVariable(s), nil)), nil
}

func (c compiler) compileUnboxed(u ast.Unboxed) (coreast.Expression, error) {
	// TODO: Handle other literals.
	switch l := u.Content().(type) {
	case ast.Number:
		return coreast.NewConstructorApplication(
			types.NewNumber(nil).CoreConstructor(),
			[]coreast.Atom{coreast.NewFloat64(l.Value())},
		), nil
	}

	panic("unreachable")
}

func (compiler) compileUnboxedLiteral(l ast.Literal) coreast.Literal {
	switch l := l.(type) {
	case ast.Number:
		return coreast.NewFloat64(l.Value())
	}

	panic("unreachable")
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
