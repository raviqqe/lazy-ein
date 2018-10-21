package types

import "llvm.org/llvm/bindings/go/llvm"

// Type is a type.
type Type interface {
	LLVMType() llvm.Type
}
