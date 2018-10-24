package types

import "llvm.org/llvm/bindings/go/llvm"

// Function is a function type.
type Function struct {
	arguments []Type
	result    Type
}

// NewFunction creates a new function type.
func NewFunction(as []Type, r Type) Function {
	return Function{as, r}
}

// LLVMType returns a LLVM type.
func (f Function) LLVMType() llvm.Type {
	return llvm.PointerType(
		llvm.StructType(
			[]llvm.Type{
				llvm.PointerType(
					llvm.FunctionType(
						f.result.LLVMType(),
						ToLLVMTypes(f.arguments),
						false,
					),
					0,
				),
			},
			false,
		),
		0,
	)
}
