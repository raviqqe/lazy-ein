package compile

import (
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

type phiGenerator struct {
	block          llvm.BasicBlock
	incomingValues []llvm.Value
	incomingBlocks []llvm.BasicBlock
}

func newPhiGenerator(b llvm.BasicBlock) *phiGenerator {
	return &phiGenerator{b, nil, nil}
}

func (g *phiGenerator) Generate(b llvm.Builder) llvm.Value {
	b.SetInsertPointAtEnd(g.block)

	v := b.CreatePHI(g.incomingValues[0].Type(), "")
	v.AddIncoming(g.incomingValues, g.incomingBlocks)

	return v
}

func (g *phiGenerator) CreateBr(b llvm.Builder, v llvm.Value) {
	b.CreateBr(g.block)
	g.incomingValues = append(g.incomingValues, v)
	g.incomingBlocks = append(g.incomingBlocks, b.GetInsertBlock())
}
