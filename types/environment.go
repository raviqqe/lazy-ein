package types

import (
	"github.com/raviqqe/stg/llir"
	"llvm.org/llvm/bindings/go/llvm"
)

// Environment is an environment type.
type Environment struct {
	size int
}

// NewEnvironment creates a new environment type.
func NewEnvironment(n int) Environment {
	return Environment{n}
}

// LLVMType returns a LLVM type.
func (e Environment) LLVMType() llvm.Type {
	return llvm.ArrayType(llvm.Int8Type(), e.size)
}

// LLVMPointerType returns a LLVM pointer type.
func (e Environment) LLVMPointerType() llvm.Type {
	return llir.PointerType(e.LLVMType())
}
