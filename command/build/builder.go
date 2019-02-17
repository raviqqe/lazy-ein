package build

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/parse"
	"llvm.org/llvm/bindings/go/llvm"
)

type builder struct {
	runtimeDirectory, moduleRootDirectory string
	objectCache                           objectCache
}

func newBuilder(runtimeDir, rootDir, cacheDir string) builder {
	return builder{runtimeDir, rootDir, newObjectCache(cacheDir, rootDir)}
}

func (b builder) Build(f string) error {
	o, err := b.BuildModule(f)

	if err != nil {
		return err
	}

	if ok, err := b.isMainModule(f); err != nil {
		return err
	} else if !ok {
		return nil
	}

	bs, err := exec.Command(
		"cc",
		o,
		b.resolveRuntimeLibrary("runtime/target/release/libio.a"),
		b.resolveRuntimeLibrary("runtime/target/release/libcore.a"),
		"-ldl",
		"-lgc",
		"-lpthread",
	).CombinedOutput()

	os.Stderr.Write(bs)

	return err
}

func (b builder) BuildModule(f string) (string, error) {
	ff, ok, err := b.objectCache.Get(f)

	if err != nil {
		return "", err
	} else if ok {
		return ff, nil
	}

	bs, err := b.buildModuleWithoutCache(f)

	if err != nil {
		return "", err
	}

	return b.objectCache.Store(f, bs)
}

func (b builder) buildModuleWithoutCache(f string) ([]byte, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return nil, err
	}

	m, err := parse.Parse(f, string(bs))

	if err != nil {
		return nil, err
	}

	mm, err := compile.Compile(m)

	if err != nil {
		return nil, err
	}

	if m.IsMainModule() {
		if err := b.renameMainFunction(mm); err != nil {
			return nil, err
		}
	}

	return b.generateModule(mm)
}

func (b builder) generateModule(m llvm.Module) ([]byte, error) {
	triple := llvm.DefaultTargetTriple()
	target, err := llvm.GetTargetFromTriple(triple)

	if err != nil {
		return nil, err
	}

	buf, err := target.CreateTargetMachine(
		triple,
		"",
		"",
		llvm.CodeGenLevelAggressive,
		llvm.RelocPIC,
		llvm.CodeModelDefault,
	).EmitToMemoryBuffer(m, llvm.ObjectFile)

	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b builder) resolveRuntimeLibrary(f string) string {
	return filepath.Join(b.runtimeDirectory, filepath.FromSlash(f))
}

func (b builder) renameMainFunction(m llvm.Module) error {
	v := m.NamedGlobal("main")

	if v.IsNil() {
		return errors.New("main function not found")
	}

	v.SetName("ein_main")

	return nil
}

func (b builder) isMainModule(f string) (bool, error) {
	bs, err := ioutil.ReadFile(f)

	if err != nil {
		return false, err
	}

	m, err := parse.Parse(f, string(bs))

	if err != nil {
		return false, err
	}

	return m.IsMainModule(), err
}
