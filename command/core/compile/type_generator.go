package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

type typeGenerator struct {
	stack      []llvm.Type
	cache      map[string]llvm.Type
	targetData llvm.TargetData
}

func newTypeGenerator(m llvm.Module) typeGenerator {
	return typeGenerator{nil, map[string]llvm.Type{}, llvm.NewTargetData(m.DataLayout())}
}

func (g typeGenerator) Generate(t types.Type) llvm.Type {
	if t, ok := g.cache[t.String()]; ok {
		return t
	}

	switch t := t.(type) {
	case types.Algebraic:
		if types.IsRecursive(t) {
			s := llvm.GlobalContext().StructCreateNamed(t.String())
			g.cache[t.String()] = s

			s.StructSetBody(g.pushType(s).generateAlgebraicBody(t), false)
			return s
		}

		return llir.StructType(g.pushDummyType().generateAlgebraicBody(t))
	case types.Boxed:
		return llir.PointerType(
			g.generateClosure(g.generateEntryFunction(nil, t.Content()), g.GenerateUnsizedPayload()),
		)
	case types.Float64:
		return llvm.DoubleType()
	case types.Function:
		if types.IsRecursive(t) {
			s := llvm.GlobalContext().StructCreateNamed(t.String())
			g.cache[t.String()] = s

			s.StructSetBody(
				g.pushType(llir.PointerType(s)).generateFunctionCloure(t).StructElementTypes(),
				false,
			)
			return llir.PointerType(s)
		}

		return llir.PointerType(g.pushDummyType().generateFunctionCloure(t))
	case types.Index:
		return g.stack[len(g.stack)-1-t.Value()]
	}

	panic("unreachable")
}

func (g typeGenerator) generateAlgebraicBody(t types.Algebraic) []llvm.Type {
	if len(t.Constructors()) == 1 {
		return g.GenerateConstructorElements(t.Constructors()[0]).StructElementTypes()
	}

	n := 0

	for _, c := range t.Constructors() {
		if m := g.GetSize(g.GenerateConstructorElements(c)); m > n {
			n = m
		}
	}

	return []llvm.Type{
		g.GenerateConstructorTag(),
		llvm.ArrayType(llvm.Int64Type(), g.bytesToWords(n)),
	}
}

func (g typeGenerator) generateFunctionCloure(t types.Function) llvm.Type {
	return g.generateClosure(
		g.generateEntryFunction(t.Arguments(), t.Result()),
		g.GenerateUnsizedPayload(),
	)
}

func (g typeGenerator) GenerateSizedClosure(l ast.LambdaDeclaration) llvm.Type {
	return g.generateClosure(g.GenerateLambdaEntryFunction(l), g.generateSizedPayload(l))
}

func (g typeGenerator) GenerateUnsizedClosure(t llvm.Type) llvm.Type {
	return g.generateClosure(t.StructElementTypes()[0].ElementType(), g.GenerateUnsizedPayload())
}

func (g typeGenerator) generateClosure(f llvm.Type, p llvm.Type) llvm.Type {
	return llir.StructType([]llvm.Type{llir.PointerType(f), p})
}

func (g typeGenerator) GenerateLambdaEntryFunction(l ast.LambdaDeclaration) llvm.Type {
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

func (g typeGenerator) generateSizedPayload(l ast.LambdaDeclaration) llvm.Type {
	n := g.GetSize(g.GenerateEnvironment(l))

	if m := g.GetSize(g.Generate(types.Unbox(l.ResultType()))); l.IsUpdatable() && m > n {
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

func (g typeGenerator) GenerateEnvironment(l ast.LambdaDeclaration) llvm.Type {
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
	return llir.WordType()
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

func (g typeGenerator) GetSize(t llvm.Type) int {
	return int(g.targetData.TypeAllocSize(t))
}

func (g typeGenerator) bytesToWords(n int) int {
	if n == 0 {
		return 0
	}

	return (n-1)/g.GetSize(llir.WordType()) + 1
}

func (g typeGenerator) pushType(t llvm.Type) typeGenerator {
	return typeGenerator{append(g.stack, t), g.cache, g.targetData}
}

func (g typeGenerator) pushDummyType() typeGenerator {
	return typeGenerator{append(g.stack, llvm.VoidType()), g.cache, g.targetData}
}
