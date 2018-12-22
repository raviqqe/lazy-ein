package desugar

import (
	"fmt"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/types"
)

func desugarUnsaturatedApplications(m ast.Module) ast.Module {
	bs := make([]ast.Bind, 0, len(m.Binds()))

	for _, b := range m.Binds() {
		bs = append(bs, desugarUnsaturatedApplicationsInBind(b))
	}

	return ast.NewModule(m.Name(), bs)
}

func desugarUnsaturatedApplicationsInBind(b ast.Bind) ast.Bind {
	e := b.Expression().ConvertExpression(
		func(e ast.Expression) ast.Expression {
			l, ok := e.(ast.Let)

			if !ok {
				return e
			}

			bs := make([]ast.Bind, 0, len(l.Binds()))

			for _, b := range l.Binds() {
				bs = append(bs, desugarUnsaturatedApplicationsInBind(b))
			}

			return ast.NewLet(bs, l.Expression())
		},
	).(ast.Expression)

	t, ok := b.Type().(types.Function)

	if !ok {
		return ast.NewBind(b.Name(), b.Type(), e)
	} else if l, ok := e.(ast.Lambda); ok && len(l.Arguments()) != t.ArgumentsCount() {
		ss, es := generateAdditionalArguments(t.ArgumentsCount() - len(l.Arguments()))
		e = ast.NewLambda(
			append(l.Arguments(), ss...),
			desugarUnsaturatedExpression(l.Expression(), es),
		)
	} else if !ok {
		ss, es := generateAdditionalArguments(t.ArgumentsCount())
		e = ast.NewLambda(ss, desugarUnsaturatedExpression(b.Expression(), es))
	}

	return ast.NewBind(b.Name(), b.Type(), e)
}

func desugarUnsaturatedExpression(e ast.Expression, as []ast.Expression) ast.Expression {
	switch e := e.(type) {
	case ast.Application:
		return ast.NewApplication(e.Function(), append(e.Arguments(), as...))
	case ast.Variable:
		return ast.NewApplication(e, as)
	case ast.Let:
		return ast.NewLet(e.Binds(), desugarUnsaturatedExpression(e.Expression(), as))
	}

	panic("unreachable")
}

func generateAdditionalArguments(n int) ([]string, []ast.Expression) {
	ss := make([]string, 0, n)
	es := make([]ast.Expression, 0, n)

	for i := 0; i < n; i++ {
		s := fmt.Sprintf("additional.argument-%v", i)
		ss = append(ss, s)
		es = append(es, ast.NewVariable(s))
	}

	return ss, es
}
