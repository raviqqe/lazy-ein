package compile

import (
	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/compile/desugar"
	coreast "github.com/raviqqe/jsonxx/command/core/ast"
	coretypes "github.com/raviqqe/jsonxx/command/core/types"
	"github.com/raviqqe/jsonxx/command/types"
)

// Compile compiles a module into a module in STG.
func Compile(m ast.Module) coreast.Module {
	m = desugar.Desugar(m)

	bs := make([]coreast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		bs = append(bs, compileBind(b))
	}

	return coreast.NewModule(m.Name(), nil, bs)
}

func compileBind(b ast.Bind) coreast.Bind {
	if len(b.Arguments()) == 0 {
		return coreast.NewBind(
			b.Name(),
			coreast.NewLambda(
				nil,
				true,
				nil,
				compileExpression(b.Expression()),
				compileType(b.Type()),
			),
		)
	}

	t := compileType(b.Type()).(coretypes.Function)
	as := make([]coreast.Argument, 0, len(b.Arguments()))

	for i, n := range b.Arguments() {
		as = append(as, coreast.NewArgument(n, t.Arguments()[i]))
	}

	return coreast.NewBind(
		b.Name(),
		coreast.NewLambda(
			nil,
			false,
			as,
			compileExpression(b.Expression()),
			t.Result(),
		),
	)
}

func compileExpression(e ast.Expression) coreast.Expression {
	switch e := e.(type) {
	case ast.Number:
		return coreast.NewFloat64(e.Value())
	case ast.Variable:
		return coreast.NewApplication(coreast.NewVariable(e.Name()), nil)
	}

	panic("unreahable")
}

func compileType(t types.Type) coretypes.Type {
	switch t := t.(type) {
	case types.Function:
		as := []coretypes.Type{}

		for {
			as = append(as, compileType(t.Argument()))
			f, ok := t.Result().(types.Function)

			if !ok {
				return coretypes.NewFunction(as, compileType(t.Result()))
			}

			t = f
		}
	case types.Unboxed:
		return compileRawType(t.Content())
	}

	return coretypes.NewBoxed(compileRawType(t))
}

func compileRawType(t types.Type) coretypes.Type {
	switch t.(type) {
	case types.Number:
		return coretypes.NewFloat64()
	}

	panic("unreachable")
}
