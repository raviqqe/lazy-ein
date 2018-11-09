package llir

import "llvm.org/llvm/bindings/go/llvm"

// CreateCall creates a common call.
func CreateCall(b llvm.Builder, f llvm.Value, as []llvm.Value) llvm.Value {
	v := b.CreateCall(f, as, "")
	v.SetInstructionCallConv(llvm.FastCallConv)
	v.SetTailCall(true)
	return v
}
