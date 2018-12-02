package codegen

import (
	"github.com/raviqqe/jsonxx/command/stg/ast"
	"github.com/raviqqe/jsonxx/command/stg/codegen/llir"
	"github.com/raviqqe/jsonxx/command/stg/codegen/names"
	"github.com/raviqqe/jsonxx/command/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

const (
	environmentArgumentName  = "environment"
	maxCodeOptimizationLevel = 3
	maxSizeOptimizationLevel = 2
)

type moduleGenerator struct {
	module          llvm.Module
	globalVariables map[string]llvm.Value
	typeGenerator   typeGenerator
}

func newModuleGenerator(m llvm.Module, ds []ast.ConstructorDefinition) (*moduleGenerator, error) {
	g := newConstructorDefinitionGenerator(m)

	for _, d := range ds {
		if err := g.GenerateUnionifyFunction(d); err != nil {
			return nil, err
		}

		if err := g.GenerateStructifyFunction(d); err != nil {
			return nil, err
		}

		g.GenerateTag(d)
	}

	return &moduleGenerator{m, map[string]llvm.Value{}, newTypeGenerator(m)}, nil
}

func (g *moduleGenerator) Generate(bs []ast.Bind) error {
	for _, b := range bs {
		g.globalVariables[b.Name()] = llvm.AddGlobal(
			g.module,
			g.typeGenerator.GenerateSizedClosure(b.Lambda()),
			b.Name(),
		)
	}

	for _, b := range bs {
		v := g.globalVariables[b.Name()]
		f, err := g.createLambda(b.Name(), b.Lambda())

		if err != nil {
			return err
		} else if err := llvm.VerifyFunction(f, llvm.AbortProcessAction); err != nil {
			return err
		}

		v.SetInitializer(
			llvm.ConstStruct(
				[]llvm.Value{f, llvm.ConstNull(v.Type().ElementType().StructElementTypes()[1])},
				false,
			),
		)
	}

	// nolint: gotype
	if err := llvm.VerifyModule(g.module, llvm.AbortProcessAction); err != nil {
		return err
	}

	g.optimize()

	return llvm.VerifyModule(g.module, llvm.AbortProcessAction)
}

func (g *moduleGenerator) createLambda(n string, l ast.Lambda) (llvm.Value, error) {
	f := llir.AddFunction(
		g.module,
		names.ToEntry(n),
		g.typeGenerator.GenerateLambdaEntryFunction(l),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	v, err := newFunctionBodyGenerator(
		b,
		g.createLogicalEnvironment(f, b, l),
		g.createLambda,
	).Generate(l.Body())

	if err != nil {
		return llvm.Value{}, err
	} else if _, ok := l.ResultType().(types.Boxed); ok && l.IsThunk() {
		// TODO: Steal child thunks in a thread-safe way.
		// TODO: Use loop to unbox children recursively.
		v = forceThunk(b, v, g.typeGenerator)
	}

	if l.IsUpdatable() {
		b.CreateStore(v, b.CreateBitCast(f.FirstParam(), llir.PointerType(v.Type()), ""))
		b.CreateStore(
			g.createUpdatedEntryFunction(n, f.Type().ElementType()),
			b.CreateGEP(
				b.CreateBitCast(f.FirstParam(), llir.PointerType(f.Type()), ""),
				[]llvm.Value{llvm.ConstIntFromString(g.typeGenerator.GenerateConstructorTag(), "-1", 10)},
				"",
			),
		)
	}

	b.CreateRet(v)

	return f, nil
}

func (g *moduleGenerator) createUpdatedEntryFunction(n string, t llvm.Type) llvm.Value {
	f := llir.AddFunction(g.module, names.ToUpdatedEntry(n), t)
	f.FirstParam().SetName(environmentArgumentName)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))
	b.CreateRet(
		b.CreateLoad(
			b.CreateBitCast(
				f.FirstParam(), llir.PointerType(f.Type().ElementType().ReturnType()),
				""),
			"",
		),
	)

	return f
}

func (g moduleGenerator) createLogicalEnvironment(f llvm.Value, b llvm.Builder, l ast.Lambda) map[string]llvm.Value {
	vs := copyVariables(g.globalVariables)

	e := b.CreateBitCast(
		f.FirstParam(),
		llir.PointerType(g.typeGenerator.GenerateEnvironment(l)),
		"",
	)

	for i, n := range l.FreeVariableNames() {
		vs[n] = b.CreateLoad(b.CreateStructGEP(e, i, ""), "")
	}

	for i, n := range append([]string{environmentArgumentName}, l.ArgumentNames()...) {
		v := f.Param(i)
		v.SetName(n)
		vs[n] = v
	}

	return vs
}

func (g moduleGenerator) optimize() {
	b := llvm.NewPassManagerBuilder()
	b.SetOptLevel(maxCodeOptimizationLevel)
	b.SetSizeLevel(maxSizeOptimizationLevel)

	p := llvm.NewPassManager()
	b.PopulateFunc(p)
	b.Populate(p)

	p.Run(g.module)
}
