package compile

import (
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

func forceThunk(b llvm.Builder, thunk llvm.Value, g typeGenerator) llvm.Value {
	return llir.CreateCall(
		b,
		b.CreateBitCast(
			b.CreateCall(
				b.GetInsertBlock().Parent().GlobalParent().NamedFunction(atomicLoadFunctionName),
				[]llvm.Value{
					b.CreateBitCast(
						b.CreateStructGEP(thunk, 0, ""),
						llir.PointerType(llir.PointerType(llvm.Int8Type())),
						"",
					),
				},
				"",
			),
			thunk.Type().ElementType().StructElementTypes()[0],
			"",
		),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(thunk, 1, ""),
				llir.PointerType(g.GenerateUnsizedPayload()),
				"",
			),
		},
	)
}
