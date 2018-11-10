package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/llir"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

type typeGenerator struct {
	module llvm.Module
}

func newTypeGenerator(m llvm.Module) typeGenerator {
	return typeGenerator{m}
}

func (g typeGenerator) GenerateSizedClosure(l ast.Lambda) llvm.Type {
	return g.generateClosure(l, g.GenerateSizedPayload(l))
}

func (g typeGenerator) GenerateUnsizedClosure(l ast.Lambda) llvm.Type {
	return g.generateClosure(l, g.GenerateUnsizedPayload())
}

func (g typeGenerator) generateClosure(l ast.Lambda, p llvm.Type) llvm.Type {
	r := l.ResultType()

	if l.IsThunk() {
		r = types.Unbox(r)
	}

	return llir.StructType(
		[]llvm.Type{
			llir.PointerType(
				llir.FunctionType(
					r.LLVMType(),
					append(
						[]llvm.Type{llir.PointerType(g.GenerateUnsizedPayload())},
						g.generateMany(l.ArgumentTypes())...,
					),
				),
			),
			p,
		},
	)
}

func (g typeGenerator) GenerateSizedPayload(l ast.Lambda) llvm.Type {
	n := g.GetSize(g.GenerateEnvironment(l))

	if m := g.GetSize(types.Unbox(l.ResultType()).LLVMType()); l.IsUpdatable() && m > n {
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

func (typeGenerator) generateMany(ts []types.Type) []llvm.Type {
	tts := make([]llvm.Type, 0, len(ts))

	for _, t := range ts {
		tts = append(tts, t.LLVMType())
	}

	return tts
}

func (g typeGenerator) GetSize(t llvm.Type) int {
	return int(llvm.NewTargetData(g.module.DataLayout()).TypeAllocSize(t))
}
