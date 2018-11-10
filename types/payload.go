package types

import (
	"github.com/raviqqe/stg/codegen/llir"
	"llvm.org/llvm/bindings/go/llvm"
)

// Payload is a payload type.
type Payload struct {
	size int
}

// NewPayload creates a new payload type.
func NewPayload(n int) Payload {
	return Payload{n}
}

// LLVMType returns a LLVM type.
func (p Payload) LLVMType() llvm.Type {
	return llvm.ArrayType(llvm.Int8Type(), p.size)
}

// LLVMPointerType returns a LLVM pointer type.
func (p Payload) LLVMPointerType() llvm.Type {
	return llir.PointerType(p.LLVMType())
}
