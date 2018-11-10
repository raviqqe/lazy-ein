package types

import (
	"github.com/raviqqe/stg/codegen/llir"
	"llvm.org/llvm/bindings/go/llvm"
)

// Boxed is a boxed type.
type Boxed struct {
	content Type
}

// NewBoxed creates a new boxed type.
func NewBoxed(t Type) Boxed {
	if _, ok := t.(Boxed); ok {
		panic("cannot box boxed types")
	}

	return Boxed{t}
}

// Content returns a content type.
func (b Boxed) Content() Type {
	return b.content
}

// LLVMType returns a LLVM type.
func (b Boxed) LLVMType() llvm.Type {
	return llir.PointerType(NewClosure(NewEnvironment(0), nil, b.content).LLVMType())
}
