package types

import "llvm.org/llvm/bindings/go/llvm"

// Unbox converts a type into its unboxed type.
func Unbox(t Type) Type {
	if t, ok := t.(Boxed); ok {
		return t.InternalType()
	}

	return t
}

// ToLLVMTypes converts types into LLVM types.
func ToLLVMTypes(ts []Type) []llvm.Type {
	ls := make([]llvm.Type, 0, len(ts))

	for _, t := range ts {
		ls = append(ls, t.LLVMType())
	}

	return ls
}
