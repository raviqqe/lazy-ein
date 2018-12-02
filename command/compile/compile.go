package compile

import (
	"github.com/raviqqe/jsonxx/command/ast"
	core "github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/types"
)

// Compile compiles a module into a module in STG.
func Compile(m ast.Module) core.Module {
	bs := make([]core.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		bs = append(bs, compileBind(b))
	}

	return core.NewModule(m.Name(), nil, bs)
}

func compileBind(b ast.Bind) core.Bind {
	return core.NewBind(
		b.Name(),
		core.NewLambda(
			nil,
			true,
			nil,
			core.NewFloat64(b.Expression().(ast.Number).Value()),
			types.NewFloat64(),
		),
	)
}
