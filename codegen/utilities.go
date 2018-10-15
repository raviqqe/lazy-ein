package codegen

import "llvm.org/llvm/bindings/go/llvm"

var voidPointerType = llvm.PointerType(llvm.VoidType(), 0)
var thunkPointerType = voidPointerType

func toEntryName(s string) string {
	return s + "-entry"
}
