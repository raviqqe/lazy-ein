package codegen

import "llvm.org/llvm/bindings/go/llvm"

var voidPointerType = llvm.PointerType(llvm.VoidType(), 0)
var thunkPointerType = llvm.PointerType(llvm.VoidType(), 0)

func toEntryName(s string) string {
	return s + "-entry"
}

func thunkPointerArrayType(n int) []llvm.Type {
	ts := make([]llvm.Type, 0, n)

	for i := 0; i < n; i++ {
		ts = append(ts, thunkPointerType)
	}

	return ts
}
