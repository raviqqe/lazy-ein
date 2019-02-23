package llir

import "github.com/llvm-mirror/llvm/bindings/go/llvm"

// PointerType creates a pointer type.
func PointerType(t llvm.Type) llvm.Type {
	return llvm.PointerType(t, 0)
}

// StructType creates a struct type.
func StructType(ts []llvm.Type) llvm.Type {
	return llvm.StructType(ts, false)
}

// FunctionType creates a function type.
func FunctionType(r llvm.Type, as []llvm.Type) llvm.Type {
	return llvm.FunctionType(r, as, false)
}

// WordType is a word type.
func WordType() llvm.Type {
	return llvm.Int64Type()
}
