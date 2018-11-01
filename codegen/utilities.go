package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

func copyVariables(vs map[string]llvm.Value) map[string]llvm.Value {
	ws := make(map[string]llvm.Value, len(vs))

	for k, v := range vs {
		ws[k] = v
	}

	return ws
}

func typeSize(m llvm.Module, t llvm.Type) int {
	return int(llvm.NewTargetData(m.DataLayout()).TypeAllocSize(t))
}

func lambdaToFreeVariablesStructType(l ast.Lambda) llvm.Type {
	return llvm.StructType(types.ToLLVMTypes(l.FreeVariableTypes()), false)
}
