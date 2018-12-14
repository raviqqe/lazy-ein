package compile

import (
	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile/desugar"
	coreast "github.com/ein-lang/ein/command/core/ast"
	corecompile "github.com/ein-lang/ein/command/core/compile"
	coretypes "github.com/ein-lang/ein/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

// Compile compiles a module into a module in STG.
func Compile(m ast.Module) (llvm.Module, error) {
	return corecompile.Compile(compileToCore(m))
}

func compileToCore(m ast.Module) coreast.Module {
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
				b.Type().ToCore(),
			),
		)
	}

	t := b.Type().ToCore().(coretypes.Function)
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
	case ast.Let:
		bs := make([]coreast.Bind, 0, len(e.Binds()))

		for _, b := range e.Binds() {
			bs = append(bs, compileBind(b))
		}

		return coreast.NewLet(bs, compileExpression(e.Expression()))
	case ast.Number:
		return coreast.NewFloat64(e.Value())
	case ast.Variable:
		return coreast.NewApplication(coreast.NewVariable(e.Name()), nil)
	}

	panic("unreahable")
}
