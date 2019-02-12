package generate

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"llvm.org/llvm/bindings/go/llvm"
)

// Executable generates an executable file.
func Executable(m llvm.Module, file, root string) error {
	if err := renameMainFunction(m); err != nil {
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
	).EmitToMemoryBuffer(m, llvm.ObjectFile)

	if err != nil {
		return err
	}

	f := strings.TrimSuffix(file, filepath.Ext(file)) + ".o"

	if err = ioutil.WriteFile(f, b.Bytes(), 0644); err != nil {
		return err
	}

	bs, err := exec.Command(
		"cc",
		resolveRuntimePath(root, "runtime/executable/libexecutable.a"),
		f,
		resolveRuntimePath(root, "runtime/runtime/target/release/libruntime.a"),
		"-ldl",
		"-lgc",
		"-lpthread",
	).CombinedOutput()

	os.Stderr.Write(bs)

	return err
}

func resolveRuntimePath(r, p string) string {
	return filepath.Join(r, filepath.FromSlash(p))
}

func renameMainFunction(m llvm.Module) error {
	v := m.NamedGlobal("main")

	if v.IsNil() {
		return errors.New("main function not found")
	}

	v.SetName("ein_main")

	return nil
}
