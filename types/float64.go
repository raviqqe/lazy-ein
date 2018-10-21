package types

import "llvm.org/llvm/bindings/go/llvm"

// Float64 is a 64-bit float type.
type Float64 struct{}

// NewFloat64 creates a new 64-bit float type.
func NewFloat64() Float64 {
	return Float64{}
}

// LLVMType returns a LLVM type.
func (f Float64) LLVMType() llvm.Type {
	return llvm.DoubleType()
}
