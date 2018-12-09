package compile

import (
	"github.com/raviqqe/jsonxx/command/ast"
	"github.com/raviqqe/jsonxx/command/compile/desugar"
	cast "github.com/raviqqe/jsonxx/command/core/ast"
	ctypes "github.com/raviqqe/jsonxx/command/core/types"
	"github.com/raviqqe/jsonxx/command/types"
)

// Compile compiles a module into a module in STG.
func Compile(m ast.Module) cast.Module {
	m = desugar.Desugar(m)

	bs := make([]cast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		bs = append(bs, compileBind(b))
	}

	return cast.NewModule(m.Name(), nil, bs)
}

func compileBind(b ast.Bind) cast.Bind {
	if len(b.Arguments()) == 0 {
		return cast.NewBind(
			b.Name(),
			cast.NewLambda(
				nil,
				true,
				nil,
				compileExpression(b.Expression()),
				compileType(b.Type()),
			),
		)
	}

	t := compileType(b.Type()).(ctypes.Function)
	as := make([]cast.Argument, 0, len(b.Arguments()))

	for i, n := range b.Arguments() {
		as = append(as, cast.NewArgument(n, t.Arguments()[i]))
	}

	return cast.NewBind(
		b.Name(),
		cast.NewLambda(
			nil,
			false,
			as,
			compileExpression(b.Expression()),
			t.Result(),
		),
	)
}

func compileExpression(e ast.Expression) cast.Expression {
	switch e := e.(type) {
	case ast.Number:
		return cast.NewFloat64(e.Value())
	case ast.Variable:
		return cast.NewApplication(cast.NewVariable(e.Name()), nil)
	}

	panic("unreahable")
}

func compileType(t types.Type) ctypes.Type {
	switch t := t.(type) {
	case types.Function:
		as := []ctypes.Type{}

		for {
			as = append(as, compileType(t.Argument()))
			f, ok := t.Result().(types.Function)

			if !ok {
				return ctypes.NewFunction(as, compileType(t.Result()))
			}

			t = f
		}
	case types.Unboxed:
		return compileRawType(t.Content())
	}

	return ctypes.NewBoxed(compileRawType(t))
}

func compileRawType(t types.Type) ctypes.Type {
	switch t.(type) {
	case types.Number:
		return ctypes.NewFloat64()
	}

	panic("unreachable")
}
