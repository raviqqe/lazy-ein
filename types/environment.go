package types

import "llvm.org/llvm/bindings/go/llvm"

// EnvironmentType is an environment type.
var EnvironmentType = llvm.ArrayType(llvm.Int8Type(), 0)

// EnvironmentPointerType is an environment pointer type.
var EnvironmentPointerType = llvm.PointerType(EnvironmentType, 0)
