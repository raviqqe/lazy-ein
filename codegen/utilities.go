package codegen

import "llvm.org/llvm/bindings/go/llvm"

var voidPointerType = llvm.PointerType(llvm.VoidType(), 0)
var thunkPointerType = llvm.PointerType(llvm.VoidType(), 0)

func toEntryName(s string) string {
	return s + "-entry"
}
