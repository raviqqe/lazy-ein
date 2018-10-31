package codegen

import "llvm.org/llvm/bindings/go/llvm"

func toEntryName(s string) string {
	return s + "-entry"
}

func toUpdatedEntryName(s string) string {
	return s + "-updated-entry"
}

// TODO: Handle duplicate names in a function.
func toInternalLambdaName(s, ss string) string {
	return s + "." + ss
}

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
