package build

import "llvm.org/llvm/bindings/go/llvm"

// nolint: gochecknoinits
func init() {
	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()
}
