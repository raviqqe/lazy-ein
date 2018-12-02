package codegen

import (
	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/codegen/llir"
	"github.com/raviqqe/jsonxx/command/core/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type typeGenerator struct {
	targetData llvm.TargetData
}

func newTypeGenerator(m llvm.Module) typeGenerator {
	return typeGenerator{llvm.NewTargetData(m.DataLayout())}
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
			g.generateClosure(g.GenerateUnsizedPayload(), g.generateEntryFunction(nil, t.Content())),
		)
	case types.Float64:
		return llvm.DoubleType()
	case types.Function:
		return llir.PointerType(
			g.generateClosure(
				g.GenerateUnsizedPayload(),
				g.generateEntryFunction(t.Arguments(), t.Result()),
			),
		)
	}

	panic("unreachable")
}

func (g typeGenerator) GenerateSizedClosure(l ast.Lambda) llvm.Type {
	return g.generateLambdaClosure(l, g.generateSizedPayload(l))
}

func (g typeGenerator) GenerateUnsizedClosure(l ast.Lambda) llvm.Type {
	return g.generateLambdaClosure(l, g.GenerateUnsizedPayload())
}

func (g typeGenerator) generateLambdaClosure(l ast.Lambda, p llvm.Type) llvm.Type {
	return g.generateClosure(p, g.GenerateLambdaEntryFunction(l))
}

func (g typeGenerator) generateClosure(p llvm.Type, f llvm.Type) llvm.Type {
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
	if g.targetData.PointerSize() < 8 {
		return llvm.Int32Type()
	}

	return llvm.Int64Type()
}

func (g typeGenerator) GenerateConstructorElements(c types.Constructor) llvm.Type {
	return llir.StructType(g.generateMany(c.Elements()))
}

func (g typeGenerator) GenerateConstructorUnionifyFunction(t types.Algebraic, i int) llvm.Type {
	return llir.FunctionType(g.Generate(t), g.generateMany(t.Constructors()[i].Elements()))
}

func (g typeGenerator) GenerateConstructorStructifyFunction(t types.Algebraic, i int) llvm.Type {
	return llir.FunctionType(
		g.GenerateConstructorElements(t.Constructors()[i]),
		[]llvm.Type{g.Generate(t)},
	)
}

func (g typeGenerator) getSize(t llvm.Type) int {
	return int(g.targetData.TypeAllocSize(t))
}

func (g typeGenerator) bytesToWords(n int) int {
	if n == 0 {
		return 0
	}

	return (n-1)/g.targetData.PointerSize() + 1
}
