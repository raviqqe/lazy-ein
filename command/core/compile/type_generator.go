package compile

import (
	"errors"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type typeGenerator struct {
	typeMap map[string]llvm.Type
	module  llvm.Module
}

func newTypeGenerator(m llvm.Module, ds []ast.TypeDefinition) (typeGenerator, error) {
	g := typeGenerator{make(map[string]llvm.Type, len(ds)), m}

	for _, d := range ds {
		if _, ok := d.Type().(types.Algebraic); ok {
			g.typeMap[d.Name()] = g.module.Context().StructCreateNamed(d.Name())
		}
	}

	for _, d := range ds {
		t, err := g.Generate(d.Type())

		if err != nil {
			return typeGenerator{}, err
		} else if _, ok := d.Type().(types.Algebraic); ok {
			g.typeMap[d.Name()].StructSetBody(t.StructElementTypes(), t.IsStructPacked())
			continue
		}

		g.typeMap[d.Name()] = t
	}

	return g, nil
}

func (g typeGenerator) Generate(t types.Type) (llvm.Type, error) {
	switch t := t.(type) {
	case types.Algebraic:
		if len(t.Constructors()) == 1 {
			return g.GenerateConstructorElements(t.Constructors()[0])
		}

		n := 0

		for _, c := range t.Constructors() {
			t, err := g.GenerateConstructorElements(c)

			if err != nil {
				return llvm.Type{}, err
			}

			if m := g.getSize(t); m > n {
				n = m
			}
		}

		return llir.StructType(
			[]llvm.Type{
				g.GenerateConstructorTag(),
				llvm.ArrayType(llvm.Int64Type(), g.bytesToWords(n)),
			},
		), nil
	case types.Boxed:
		tt, err := g.generateEntryFunction(nil, t.Content())

		if err != nil {
			return llvm.Type{}, err
		}

		return llir.PointerType(g.generateClosure(tt, g.GenerateUnsizedPayload())), nil
	case types.Float64:
		return llvm.DoubleType(), nil
	case types.Function:
		tt, err := g.generateEntryFunction(t.Arguments(), t.Result())

		if err != nil {
			return llvm.Type{}, err
		}

		return llir.PointerType(g.generateClosure(tt, g.GenerateUnsizedPayload())), nil
	case types.Named:
		if t, ok := g.typeMap[t.Name()]; ok {
			return t, nil
		}

		return llvm.Type{}, errors.New("type " + t.Name() + " is undefined")
	}

	panic("unreachable")
}

func (g typeGenerator) GenerateSizedClosure(l ast.Lambda) (llvm.Type, error) {
	e, err := g.GenerateLambdaEntryFunction(l)

	if err != nil {
		return llvm.Type{}, err
	}

	p, err := g.generateSizedPayload(l)

	if err != nil {
		return llvm.Type{}, err
	}

	return g.generateClosure(e, p), nil
}

func (g typeGenerator) GenerateUnsizedClosure(t llvm.Type) llvm.Type {
	return g.generateClosure(t.StructElementTypes()[0].ElementType(), g.GenerateUnsizedPayload())
}

func (g typeGenerator) generateClosure(f llvm.Type, p llvm.Type) llvm.Type {
	return llir.StructType([]llvm.Type{llir.PointerType(f), p})
}

func (g typeGenerator) GenerateLambdaEntryFunction(l ast.Lambda) (llvm.Type, error) {
	r := l.ResultType()

	if l.IsThunk() {
		r = types.Unbox(r)
	}

	return g.generateEntryFunction(l.ArgumentTypes(), r)
}

func (g typeGenerator) generateEntryFunction(as []types.Type, r types.Type) (llvm.Type, error) {
	t, err := g.Generate(r)

	if err != nil {
		return llvm.Type{}, err
	}

	ts, err := g.generateMany(as)

	if err != nil {
		return llvm.Type{}, err
	}

	return llir.FunctionType(
		t,
		append(
			[]llvm.Type{llir.PointerType(g.GenerateUnsizedPayload())},
			ts...,
		),
	), nil
}

func (g typeGenerator) generateSizedPayload(l ast.Lambda) (llvm.Type, error) {
	t, err := g.GenerateEnvironment(l)

	if err != nil {
		return llvm.Type{}, err
	}

	n := g.getSize(t)

	if t, err := g.Generate(types.Unbox(l.ResultType())); err != nil {
		return llvm.Type{}, err
	} else if m := g.getSize(t); l.IsUpdatable() && m > n {
		n = m
	}

	return g.generatePayload(n), nil
}

func (g typeGenerator) GenerateUnsizedPayload() llvm.Type {
	return g.generatePayload(0)
}

func (g typeGenerator) generatePayload(n int) llvm.Type {
	return llvm.ArrayType(llvm.Int8Type(), n)
}

func (g typeGenerator) GenerateEnvironment(l ast.Lambda) (llvm.Type, error) {
	ts, err := g.generateMany(l.FreeVariableTypes())

	if err != nil {
		return llvm.Type{}, err
	}

	return llir.StructType(ts), nil
}

func (g typeGenerator) generateMany(ts []types.Type) ([]llvm.Type, error) {
	tts := make([]llvm.Type, 0, len(ts))

	for _, t := range ts {
		t, err := g.Generate(t)

		if err != nil {
			return nil, err
		}

		tts = append(tts, t)
	}

	return tts, nil
}

func (g typeGenerator) GenerateConstructorTag() llvm.Type {
	if g.targetData().PointerSize() < 8 {
		return llvm.Int32Type()
	}

	return llvm.Int64Type()
}

func (g typeGenerator) GenerateConstructorElements(c types.Constructor) (llvm.Type, error) {
	ts, err := g.generateMany(c.Elements())

	if err != nil {
		return llvm.Type{}, err
	}

	return llir.StructType(ts), nil
}

func (g typeGenerator) GenerateConstructorUnionifyFunction(
	a types.Algebraic,
	c types.Constructor,
) (llvm.Type, error) {
	t, err := g.Generate(a)

	if err != nil {
		return llvm.Type{}, err
	}

	ts, err := g.generateMany(c.Elements())

	if err != nil {
		return llvm.Type{}, err
	}

	return llir.FunctionType(t, ts), nil
}

func (g typeGenerator) GenerateConstructorStructifyFunction(
	a types.Algebraic,
	c types.Constructor,
) (llvm.Type, error) {
	t, err := g.GenerateConstructorElements(c)

	if err != nil {
		return llvm.Type{}, err
	}

	tt, err := g.Generate(a)

	if err != nil {
		return llvm.Type{}, err
	}

	return llir.FunctionType(t, []llvm.Type{tt}), nil
}

func (g typeGenerator) getSize(t llvm.Type) int {
	return int(g.targetData().TypeAllocSize(t))
}

func (g typeGenerator) bytesToWords(n int) int {
	if n == 0 {
		return 0
	}

	return (n-1)/g.targetData().PointerSize() + 1
}

func (g typeGenerator) targetData() llvm.TargetData {
	return llvm.NewTargetData(g.module.DataLayout())
}
