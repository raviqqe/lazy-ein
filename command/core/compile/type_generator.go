package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type typeGenerator struct {
	typeMap map[string]llvm.Type
	module  llvm.Module
}

func newTypeGenerator(m llvm.Module, ds []ast.TypeDefinition) typeGenerator {
	g := typeGenerator{make(map[string]llvm.Type, len(ds)), m}

	for _, d := range ds {
		if _, ok := d.Type().(types.Algebraic); ok {
			g.typeMap[d.Name()] = g.module.Context().StructCreateNamed(d.Name())
		}
	}

	for _, d := range ds {
		t := g.Generate(d.Type())

		if _, ok := d.Type().(types.Algebraic); !ok {
			g.typeMap[d.Name()].StructSetBody(t.StructElementTypes(), t.IsStructPacked())
			continue
		}

		g.typeMap[d.Name()] = t
	}

	return g
}

func (g typeGenerator) Generate(t types.Type) llvm.Type {
	switch t := t.(type) {
	case types.Algebraic:
		if len(t.Constructors()) == 1 {
			return g.GenerateConstructorElements(t.Constructors()[0])
		}

		n := 0

		for _, c := range t.Constructors() {
			if m := g.getSize(g.GenerateConstructorElements(c)); m > n {
				n = m
			}
		}

		return llir.StructType(
			[]llvm.Type{
				g.GenerateConstructorTag(),
				llvm.ArrayType(llvm.Int64Type(), g.bytesToWords(n)),
			},
		)
	case types.Boxed:
		return llir.PointerType(
			g.generateClosure(g.generateEntryFunction(nil, t.Content()), g.GenerateUnsizedPayload()),
		)
	case types.Float64:
		return llvm.DoubleType()
	case types.Function:
		return llir.PointerType(
			g.generateClosure(
				g.generateEntryFunction(t.Arguments(), t.Result()),
				g.GenerateUnsizedPayload(),
			),
		)
	}

	panic("unreachable")
}

func (g typeGenerator) GenerateSizedClosure(l ast.Lambda) llvm.Type {
	return g.generateClosure(g.GenerateLambdaEntryFunction(l), g.generateSizedPayload(l))
}

func (g typeGenerator) GenerateUnsizedClosure(t llvm.Type) llvm.Type {
	return g.generateClosure(t.StructElementTypes()[0].ElementType(), g.GenerateUnsizedPayload())
}

func (g typeGenerator) generateClosure(f llvm.Type, p llvm.Type) llvm.Type {
	return llir.StructType([]llvm.Type{llir.PointerType(f), p})
}

func (g typeGenerator) GenerateLambdaEntryFunction(l ast.Lambda) llvm.Type {
	r := l.ResultType()

	if l.IsThunk() {
		r = types.Unbox(r)
	}

	return g.generateEntryFunction(l.ArgumentTypes(), r)
}

func (g typeGenerator) generateEntryFunction(as []types.Type, r types.Type) llvm.Type {
	return llir.FunctionType(
		g.Generate(r),
		append(
			[]llvm.Type{llir.PointerType(g.GenerateUnsizedPayload())},
			g.generateMany(as)...,
		),
	)
}

func (g typeGenerator) generateSizedPayload(l ast.Lambda) llvm.Type {
	n := g.getSize(g.GenerateEnvironment(l))

	if m := g.getSize(g.Generate(types.Unbox(l.ResultType()))); l.IsUpdatable() && m > n {
		n = m
	}

	return g.generatePayload(n)
}

func (g typeGenerator) GenerateUnsizedPayload() llvm.Type {
	return g.generatePayload(0)
}

func (g typeGenerator) generatePayload(n int) llvm.Type {
	return llvm.ArrayType(llvm.Int8Type(), n)
}

func (g typeGenerator) GenerateEnvironment(l ast.Lambda) llvm.Type {
	return llir.StructType(g.generateMany(l.FreeVariableTypes()))
}

func (g typeGenerator) generateMany(ts []types.Type) []llvm.Type {
	tts := make([]llvm.Type, 0, len(ts))

	for _, t := range ts {
		tts = append(tts, g.Generate(t))
	}

	return tts
}

func (g typeGenerator) GenerateConstructorTag() llvm.Type {
	if g.targetData().PointerSize() < 8 {
		return llvm.Int32Type()
	}

	return llvm.Int64Type()
}

func (g typeGenerator) GenerateConstructorElements(c types.Constructor) llvm.Type {
	return llir.StructType(g.generateMany(c.Elements()))
}

func (g typeGenerator) GenerateConstructorUnionifyFunction(
	a types.Algebraic,
	c types.Constructor,
) llvm.Type {
	return llir.FunctionType(g.Generate(a), g.generateMany(c.Elements()))
}

func (g typeGenerator) GenerateConstructorStructifyFunction(
	a types.Algebraic,
	c types.Constructor,
) llvm.Type {
	return llir.FunctionType(
		g.GenerateConstructorElements(c),
		[]llvm.Type{g.Generate(a)},
	)
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
