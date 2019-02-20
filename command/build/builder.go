package build

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ein-lang/ein/command/ast"
	"github.com/ein-lang/ein/command/compile"
	"github.com/ein-lang/ein/command/compile/metadata"
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
	ss, _, err := b.buildModule(f)

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
		append(
			ss,
			b.resolveRuntimeLibrary("runtime/target/release/libio.a"),
			b.resolveRuntimeLibrary("runtime/target/release/libcore.a"),
			"-ldl",
			"-lgc",
			"-lpthread",
		)...,
	).CombinedOutput()

	os.Stderr.Write(bs)

	return err
}

func (b builder) buildModule(f string) ([]string, metadata.Module, error) {
	m, err := parse.Parse(f, b.moduleRootDirectory)

	if err != nil {
		return nil, metadata.Module{}, err
	}

	ss, mds, err := b.buildSubmodules(m)

	if err != nil {
		return nil, metadata.Module{}, err
	}

	s, ok, err := b.objectCache.Get(f)

	if err != nil {
		return nil, metadata.Module{}, err
	} else if ok {
		return append(ss, s), metadata.NewModule(m), nil
	}

	mm, err := compile.Compile(m, mds)

	if err != nil {
		return nil, metadata.Module{}, err
	}

	bs, err := b.generateModule(mm)

	if err != nil {
		return nil, metadata.Module{}, err
	}

	s, err = b.objectCache.Store(f, bs)

	if err != nil {
		return nil, metadata.Module{}, err
	}

	return append(ss, s), metadata.NewModule(m), nil
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

func (b builder) buildSubmodules(m ast.Module) ([]string, []metadata.Module, error) {
	ss := make([]string, 0, len(m.Imports()))
	mds := make([]metadata.Module, 0, len(m.Imports()))

	for _, i := range m.Imports() {
		sss, md, err := b.buildModule(i.Name().ToPath(b.moduleRootDirectory))

		if err != nil {
			return nil, nil, err
		}

		ss = append(ss, sss...)
		mds = append(mds, md)
	}

	return ss, mds, nil
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
