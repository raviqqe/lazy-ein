package codegen

import (
	"github.com/raviqqe/stg/codegen/llir"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func copyVariables(vs map[string]llvm.Value) map[string]llvm.Value {
	ws := make(map[string]llvm.Value, len(vs))

	for k, v := range vs {
		ws[k] = v
	}

	return ws
}

func forceThunk(b llvm.Builder, v llvm.Value) llvm.Value {
	return llir.CreateCall(
		b,
		b.CreateLoad(b.CreateStructGEP(v, 0, ""), ""),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(v, 1, ""),
				types.NewPayload(0).LLVMPointerType(),
				"",
			),
		},
	)
}
