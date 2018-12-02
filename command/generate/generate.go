package generate

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/raviqqe/jsonxx/command/core/ast"
	"github.com/raviqqe/jsonxx/command/core/compile"
	"llvm.org/llvm/bindings/go/llvm"
)

// Executable generates an executable file.
func Executable(f string, m ast.Module) error {
	mm, err := compile.Compile(m)

	if err != nil {
		return err
	}

	s := llvm.DefaultTargetTriple()
	t, err := llvm.GetTargetFromTriple(s)

	if err != nil {
		return err
	}

	b, err := t.CreateTargetMachine(
		s,
		"",
		"",
		llvm.CodeGenLevelAggressive,
		llvm.RelocDefault,
		llvm.CodeModelDefault,
	).EmitToMemoryBuffer(mm, llvm.ObjectFile)

	if err != nil {
		return err
	}

	// TODO: Generate executable files.
	return ioutil.WriteFile(
		strings.TrimSuffix(f, filepath.Ext(f))+".o",
		b.Bytes(),
		0644,
	)
}
