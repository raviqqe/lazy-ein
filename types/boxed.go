package types

import "llvm.org/llvm/bindings/go/llvm"

// Boxed is a boxed type.
type Boxed struct {
	internalType Type
}

// NewBoxed creates a new boxed type.
func NewBoxed(t Type) Boxed {
	return Boxed{t}
}

// LLVMType returns a LLVM type.
func (b Boxed) LLVMType() llvm.Type {
	e := NewEnvironment(0)

	return llvm.PointerType(
		llvm.StructType(
			[]llvm.Type{
				llvm.PointerType(
					llvm.FunctionType(
						b.internalType.LLVMType(),
						[]llvm.Type{e.LLVMPointerType()},
						false,
					),
					0,
				),
				e.LLVMType(),
			},
			false,
		),
		0,
	)
}
