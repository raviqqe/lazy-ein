package compile

import (
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

func forceThunk(b llvm.Builder, v llvm.Value, g typeGenerator) llvm.Value {
	f := llir.CreateCall(
		b,
		b.GetInsertBlock().Parent().GlobalParent().NamedFunction(atomicLoadFunctionName),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(v, 0, ""),
				llir.PointerType(llir.PointerType(llvm.Int8Type())),
				"",
			),
		},
	)

	// TODO: Throw thunks into a black hole if f is null.

	return llir.CreateCall(
		b,
		b.CreateBitCast(f, v.Type().ElementType().StructElementTypes()[0], ""),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(v, 1, ""),
				llir.PointerType(g.GenerateUnsizedPayload()),
				"",
			),
		},
	)
}
