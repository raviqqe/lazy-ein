package compile

import (
	"github.com/raviqqe/jsonxx/command/ast"
	cast "github.com/raviqqe/jsonxx/command/core/ast"
	ctypes "github.com/raviqqe/jsonxx/command/core/types"
	"github.com/raviqqe/jsonxx/command/types"
)

// Compile compiles a module into a module in STG.
func Compile(m ast.Module) cast.Module {
	bs := make([]cast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		bs = append(bs, compileBind(b))
	}

	return cast.NewModule(m.Name(), nil, bs)
}

func compileBind(b ast.Bind) cast.Bind {
	return cast.NewBind(
		b.Name(),
		cast.NewLambda(
			nil,
			true,
			nil,
			cast.NewFloat64(b.Expression().(ast.Number).Value()),
			compileType(b.Type()),
		),
	)
}

func compileType(t types.Type) ctypes.Type {
	switch t.(type) {
	case types.Number:
		return ctypes.NewFloat64()
	}

	panic("unreahable")
}
