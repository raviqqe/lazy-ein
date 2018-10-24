package types

import "llvm.org/llvm/bindings/go/llvm"

// EnvironmentPointerType is an environment pointer type.
var EnvironmentPointerType = llvm.PointerType(NewBoxed(NewFloat64()).LLVMType(), 0)
