package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/llir"
	"github.com/raviqqe/stg/codegen/names"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

const environmentArgumentName = "environment"

type moduleGenerator struct {
	module          llvm.Module
	globalVariables map[string]llvm.Value
	typeGenerator   typeGenerator
}

func newModuleGenerator(m llvm.Module) *moduleGenerator {
	return &moduleGenerator{m, map[string]llvm.Value{}, newTypeGenerator(m)}
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

	return llvm.VerifyModule(g.module, llvm.AbortProcessAction)
}

func (g *moduleGenerator) createLambda(n string, l ast.Lambda) (llvm.Value, error) {
	t := g.typeGenerator.Generate(l.ResultType())

	if l.IsThunk() {
		t = g.typeGenerator.Generate(types.Unbox(l.ResultType()))
	}

	f := llir.AddFunction(
		g.module,
		names.ToEntry(n),
		llir.FunctionType(
			t,
			append(
				[]llvm.Type{llir.PointerType(g.typeGenerator.GenerateUnsizedPayload())},
				g.typeGenerator.generateMany(l.ArgumentTypes())...,
			),
		),
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
			g.createUpdatedEntryFunction(n, t),
			g.payloadToEntryFunctionPointer(b, f.FirstParam(), t),
		)
	}

	b.CreateRet(v)

	return f, nil
}

func (g *moduleGenerator) createUpdatedEntryFunction(n string, t llvm.Type) llvm.Value {
	f := llir.AddFunction(
		g.module,
		names.ToUpdatedEntry(n),
		llir.FunctionType(t, []llvm.Type{llir.PointerType(g.typeGenerator.GenerateUnsizedPayload())}),
	)
	f.FirstParam().SetName(environmentArgumentName)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))
	b.CreateRet(b.CreateLoad(b.CreateBitCast(f.FirstParam(), llir.PointerType(t), ""), ""))

	return f
}

func (g *moduleGenerator) payloadToEntryFunctionPointer(
	b llvm.Builder, v llvm.Value, t llvm.Type,
) llvm.Value {
	return b.CreateGEP(
		b.CreateBitCast(
			v,
			llir.PointerType(
				llir.PointerType(
					llir.FunctionType(t, []llvm.Type{llir.PointerType(g.typeGenerator.GenerateUnsizedPayload())}),
				),
			),
			"",
		),
		[]llvm.Value{llvm.ConstIntFromString(llvm.Int32Type(), "-1", 10)},
		"",
	)
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
