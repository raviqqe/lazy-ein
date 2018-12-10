package generate

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
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

	if err := renameMainFunction(mm); err != nil {
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
		llvm.RelocPIC,
		llvm.CodeModelDefault,
	).EmitToMemoryBuffer(mm, llvm.ObjectFile)

	if err != nil {
		return err
	}

	f = strings.TrimSuffix(f, filepath.Ext(f)) + ".o"

	if err = ioutil.WriteFile(f, b.Bytes(), 0644); err != nil {
		return err
	}

	r, err := getRuntimeRoot()

	if err != nil {
		return err
	}

	return exec.Command(
		"cc",
		filepath.Join(r, filepath.FromSlash("runtime/executable/libexecutable.a")),
		f,
		filepath.Join(r, filepath.FromSlash("runtime/io/target/release/libio.a")),
		"-lpthread",
		"-ldl",
	).Run()
}

func getRuntimeRoot() (string, error) {
	s := os.Getenv("JSONXX_ROOT")

	if s == "" {
		return "", errors.New("JSONXX_ROOT environment variable not set")
	}

	return s, nil
}

func renameMainFunction(m llvm.Module) error {
	v := m.NamedGlobal("main")

	if v.IsNil() {
		return errors.New("main function not found")
	}

	v.SetName("jsonxx_main")

	return nil
}
