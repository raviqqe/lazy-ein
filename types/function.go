package types

import (
	"github.com/raviqqe/stg/codegen/llir"
	"llvm.org/llvm/bindings/go/llvm"
)

// Function is a function type.
type Function struct {
	arguments []Type
	result    Type
}

// NewFunction creates a new function type.
func NewFunction(as []Type, r Type) Function {
	return Function{as, r}
}

// Arguments returns arguments.
func (f Function) Arguments() []Type {
	return f.arguments
}

// Result returns arguments.
func (f Function) Result() Type {
	return f.result
}

// LLVMType returns a LLVM type.
func (f Function) LLVMType() llvm.Type {
	return llir.PointerType(NewClosure(NewPayload(0).LLVMType(), f.arguments, f.result).LLVMType())
}
