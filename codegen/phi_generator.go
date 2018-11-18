package codegen

import (
	"github.com/raviqqe/stg/ast"
	"llvm.org/llvm/bindings/go/llvm"
)

type phiGenerator struct {
	block          llvm.BasicBlock
	incomingValues []llvm.Value
	incomingBlocks []llvm.BasicBlock
}

func newPhiGenerator(b llvm.BasicBlock, c ast.Case) *phiGenerator {
	return &phiGenerator{
		b,
		make([]llvm.Value, 0, len(c.Alternatives())+1),
		make([]llvm.BasicBlock, 0, len(c.Alternatives())+1),
	}
}

func (g *phiGenerator) Generate(b llvm.Builder) llvm.Value {
	b.SetInsertPointAtEnd(g.block)

	v := b.CreatePHI(g.incomingValues[0].Type(), "")
	v.AddIncoming(g.incomingValues, g.incomingBlocks)

	return v
}

func (g *phiGenerator) AddIncoming(v llvm.Value, b llvm.BasicBlock) {
	g.incomingValues = append(g.incomingValues, v)
	g.incomingBlocks = append(g.incomingBlocks, b)
}

func (g phiGenerator) Block() llvm.BasicBlock {
	return g.block
}
