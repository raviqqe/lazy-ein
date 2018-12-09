package compile

import (
	"github.com/raviqqe/jsonxx/command/core/compile/llir"
	"llvm.org/llvm/bindings/go/llvm"
)

func forceThunk(b llvm.Builder, v llvm.Value, g typeGenerator) llvm.Value {
	return llir.CreateCall(
		b,
		b.CreateLoad(b.CreateStructGEP(v, 0, ""), ""),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(v, 1, ""),
				llir.PointerType(g.GenerateUnsizedPayload()),
				"",
			),
		},
	)
}
