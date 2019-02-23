package llir

import "github.com/llvm-mirror/llvm/bindings/go/llvm"

// AddFunction adds a function to a module.
func AddFunction(m llvm.Module, s string, t llvm.Type) llvm.Value {
	f := llvm.AddFunction(m, s, t)
	f.SetFunctionCallConv(llvm.FastCallConv)
	return f
}
