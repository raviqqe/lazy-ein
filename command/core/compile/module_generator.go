package compile

import (
	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

const environmentArgumentName = "environment"

type moduleGenerator struct {
	module          llvm.Module
	globalVariables map[string]llvm.Value
	typeGenerator   typeGenerator
}

func newModuleGenerator() *moduleGenerator {
	return &moduleGenerator{llvm.Module{}, map[string]llvm.Value{}, typeGenerator{}}
}

func (g *moduleGenerator) initialize(m ast.Module) error {
	g.module = llvm.NewModule("")
	g.typeGenerator = newTypeGenerator(g.module)

	llvm.AddFunction(
		g.module,
		allocFunctionName,
		llvm.FunctionType(llir.PointerType(llvm.Int8Type()), []llvm.Type{llir.WordType()}, false),
	)
	llvm.AddFunction(
		g.module,
		blackHoleFunctionName,
		llvm.FunctionType(
			llir.PointerType(llvm.VoidType()),
			[]llvm.Type{llir.PointerType(llvm.Int8Type())},
			false,
		),
	)
	llvm.AddFunction(
		g.module,
		panicFunctionName,
		llvm.FunctionType(llvm.VoidType(), nil, false),
	)

	for _, d := range m.Declarations() {
		g.globalVariables[d.Name()] = llvm.AddGlobal(
			g.module,
			g.typeGenerator.GenerateSizedClosure(d.Lambda()),
			d.Name(),
		)
	}

	cg := newConstructorGenerator(g.module, g.typeGenerator)

	for _, t := range m.Types() {
		if t, ok := t.(types.Algebraic); ok {
			if err := cg.Generate(t); err != nil {
				return err
			}
		}
	}

	return nil
}

func (g *moduleGenerator) Generate(m ast.Module) (llvm.Module, error) {
	if err := g.initialize(m); err != nil {
		return llvm.Module{}, err
	}

	for _, b := range m.Binds() {
		g.globalVariables[b.Name()] = llvm.AddGlobal(
			g.module,
			g.typeGenerator.GenerateSizedClosure(b.Lambda().ToDeclaration()),
			b.Name(),
		)
	}

	for _, b := range m.Binds() {
		v := g.globalVariables[b.Name()]
		f, err := g.createLambda(b.Name(), b.Lambda())

		if err != nil {
			return llvm.Module{}, err
		} else if err := llvm.VerifyFunction(f, llvm.AbortProcessAction); err != nil {
			return llvm.Module{}, err
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
		return llvm.Module{}, err
	}

	return g.module, nil
}

func (g *moduleGenerator) createLambda(n string, l ast.Lambda) (llvm.Value, error) {
	f := llir.AddFunction(
		g.module,
		names.ToEntry(n),
		g.typeGenerator.GenerateLambdaEntryFunction(l.ToDeclaration()),
	)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))

	v, err := newFunctionBodyGenerator(
		b,
		g.createLogicalEnvironment(f, b, l),
		g.createLambda,
		g.typeGenerator,
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
	f := llir.AddFunction(g.module, names.ToNormalFormEntry(n), t)
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
	vs := make(map[string]llvm.Value, len(g.globalVariables))

	for k, v := range g.globalVariables {
		if v.Type().ElementType().StructElementTypes()[1].ArrayLength() != 0 {
			v = b.CreateBitCast(
				v,
				llir.PointerType(g.typeGenerator.GenerateUnsizedClosure(v.Type().ElementType())),
				"",
			)
		}

		vs[k] = v
	}

	e := b.CreateBitCast(
		f.FirstParam(),
		llir.PointerType(g.typeGenerator.GenerateEnvironment(l.ToDeclaration())),
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
