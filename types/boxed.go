package types

import "llvm.org/llvm/bindings/go/llvm"

// Boxed is a boxed type.
type Boxed struct {
	internalType Type
}

// NewBoxed creates a new boxed type.
func NewBoxed(t Type) Boxed {
	if _, ok := t.(Boxed); ok {
		panic("cannot box boxed types")
	}

	return Boxed{t}
}

// InternalType returns an internal type.
func (b Boxed) InternalType() Type {
	return b.internalType
}

// LLVMType returns a LLVM type.
func (b Boxed) LLVMType() llvm.Type {
	return llvm.PointerType(NewClosure(NewEnvironment(0), nil, b.internalType).LLVMType(), 0)
}
