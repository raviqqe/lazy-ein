package build

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/compile/metadata"
	"github.com/ein-lang/ein/command/parse"
	"llvm.org/llvm/bindings/go/llvm"
)

const (
	maxCodeOptimizationLevel = 3
	maxSizeOptimizationLevel = 2
)

type builder struct {
	runtimeDirectory, moduleRootDirectory string
	objectCache                           objectCache
}

func newBuilder(runtimeDir, rootDir, cacheDir string) builder {
	return builder{runtimeDir, rootDir, newObjectCache(cacheDir, rootDir)}
}

func (b builder) Build(fname string) error {
	m, _, err := b.buildModule(fname)

	if err != nil {
		return err
	}

	if ok, err := b.isMainModule(fname); err != nil {
		return err
	} else if !ok {
		return nil
	}

	bs, err := b.generateModule(m)

	if err != nil {
		return err
	}

	f, err := ioutil.TempFile("", "")

	if err != nil {
		return err
	}

	if _, err := f.Write(bs); err != nil {
		return err
	}

	defer os.Remove(f.Name())

	bs, err = exec.Command(
		"cc",
		f.Name(),
		b.resolveRuntimeLibrary("runtime/target/release/libio.a"),
		b.resolveRuntimeLibrary("runtime/target/release/libcore.a"),
		"-ldl",
		"-lgc",
		"-lpthread",
	).CombinedOutput()

	os.Stderr.Write(bs)

	return err
}

func (b builder) buildModule(f string) (llvm.Module, metadata.Module, error) {
	m, err := parse.Parse(f, b.moduleRootDirectory)

	if err != nil {
		return llvm.Module{}, metadata.Module{}, err
	}

	ms, mds, err := b.buildSubmodules(m)

	if err != nil {
		return llvm.Module{}, metadata.Module{}, err
	}

	mm, ok, err := b.objectCache.Get(f)

	if err != nil {
		return llvm.Module{}, metadata.Module{}, err
	} else if ok {
		return mm, metadata.NewModule(m), nil
	}

	mm, err = compile.Compile(m, mds)

	if err != nil {
		return llvm.Module{}, metadata.Module{}, err
	}

	if err := b.linkLLVMModules(mm, ms); err != nil {
		return llvm.Module{}, metadata.Module{}, err
	}

	if err = b.objectCache.Store(f, mm); err != nil {
		return llvm.Module{}, metadata.Module{}, err
	}

	return mm, metadata.NewModule(m), nil
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

func (b builder) buildSubmodules(m ast.Module) ([]llvm.Module, []metadata.Module, error) {
	ms := make([]llvm.Module, 0, len(m.Imports()))
	mds := make([]metadata.Module, 0, len(m.Imports()))

	for _, i := range m.Imports() {
		m, md, err := b.buildModule(i.Name().ToPath(b.moduleRootDirectory))

		if err != nil {
			return nil, nil, err
		}

		ms = append(ms, m)
		mds = append(mds, md)
	}

	return ms, mds, nil
}

func (b builder) resolveRuntimeLibrary(f string) string {
	return filepath.Join(b.runtimeDirectory, filepath.FromSlash(f))
}

func (b builder) isMainModule(f string) (bool, error) {
	m, err := parse.Parse(f, b.moduleRootDirectory)

	if err != nil {
		return false, err
	}

	return m.IsMainModule(), err
}

func (b builder) linkLLVMModules(m llvm.Module, ms []llvm.Module) error {
	for _, mm := range ms {
		if err := llvm.LinkModules(m, mm); err != nil {
			return err
		}
	}

	b.optimize(m)

	return llvm.VerifyModule(m, llvm.AbortProcessAction)
}

func (builder) optimize(m llvm.Module) {
	b := llvm.NewPassManagerBuilder()
	b.SetOptLevel(maxCodeOptimizationLevel)
	b.SetSizeLevel(maxSizeOptimizationLevel)

	p := llvm.NewPassManager()
	b.PopulateFunc(p)
	b.Populate(p)

	p.Run(m)
}
