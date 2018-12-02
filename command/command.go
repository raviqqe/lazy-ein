package command

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/raviqqe/jsonxx/command/compile"
	core "github.com/raviqqe/jsonxx/command/core/compile"
	"github.com/raviqqe/jsonxx/command/parse"
	"llvm.org/llvm/bindings/go/llvm"
)

// Command runs a compiler command with command-line arguments.
func Command(ss []string) error {
	as, err := getArguments(ss)

	if err != nil {
		return err
	}

	bs, err := ioutil.ReadFile(as.Filename)

	if err != nil {
		return err
	}

	m, err := parse.Parse(as.Filename, string(bs))

	if err != nil {
		return err
	}

	mm, err := core.Compile(compile.Compile(m))

	if err != nil {
		return err
	}

	llvm.InitializeAllTargets()
	llvm.InitializeAllTargetInfos()
	llvm.InitializeAllTargetMCs()
	llvm.InitializeAllAsmParsers()
	llvm.InitializeAllAsmPrinters()

	t, err := llvm.GetTargetFromTriple(llvm.DefaultTargetTriple())

	if err != nil {
		return err
	}

	b, err := t.CreateTargetMachine(
		mm.Target(),
		"",
		"",
		llvm.CodeGenLevelDefault,
		llvm.RelocDefault,
		llvm.CodeModelDefault,
	).EmitToMemoryBuffer(mm, llvm.ObjectFile)

	if err != nil {
		return err
	}

	return ioutil.WriteFile(
		strings.TrimSuffix(as.Filename, filepath.Ext(as.Filename))+".o",
		b.Bytes(),
		0644,
	)
}
