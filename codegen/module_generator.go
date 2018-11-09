package codegen

import (
	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/names"
	"github.com/raviqqe/stg/llir"
	"github.com/raviqqe/stg/types"
	"llvm.org/llvm/bindings/go/llvm"
)

const environmentArgumentName = "environment"

type moduleGenerator struct {
	module          llvm.Module
	globalVariables map[string]llvm.Value
}

func newModuleGenerator(m llvm.Module) *moduleGenerator {
	return &moduleGenerator{m, map[string]llvm.Value{}}
}

func (g *moduleGenerator) Generate(bs []ast.Bind) error {
	for _, b := range bs {
		g.createClosure(b.Name(), b.Lambda())
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
	t := l.ResultType().LLVMType()

	if l.IsThunk() {
		t = types.Unbox(l.ResultType()).LLVMType()
	}

	f := llir.AddFunction(
		g.module,
		names.ToEntry(n),
		llir.FunctionType(
			t,
			append(
				[]llvm.Type{types.NewEnvironment(0).LLVMPointerType()},
				types.ToLLVMTypes(l.ArgumentTypes())...,
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
		v = forceThunk(b, v)
	}

	if l.IsUpdatable() {
		b.CreateStore(v, b.CreateBitCast(f.FirstParam(), llir.PointerType(v.Type()), ""))
		b.CreateStore(
			g.createUpdatedEntryFunction(n, t),
			g.environmentToEntryFunctionPointer(b, f.FirstParam(), t),
		)
	}

	b.CreateRet(v)

	return f, nil
}

func (g *moduleGenerator) createClosure(n string, l ast.Lambda) {
	r := l.ResultType()

	if l.IsThunk() {
		r = types.Unbox(r)
	}

	g.globalVariables[n] = llvm.AddGlobal(
		g.module,
		types.NewClosure(g.lambdaToEnvironment(l), l.ArgumentTypes(), r).LLVMType(),
		n,
	)
}

func (g *moduleGenerator) createUpdatedEntryFunction(n string, t llvm.Type) llvm.Value {
	f := llir.AddFunction(
		g.module,
		names.ToUpdatedEntry(n),
		llir.FunctionType(t, []llvm.Type{types.NewEnvironment(0).LLVMPointerType()}),
	)
	f.FirstParam().SetName(environmentArgumentName)

	b := llvm.NewBuilder()
	b.SetInsertPointAtEnd(llvm.AddBasicBlock(f, ""))
	b.CreateRet(b.CreateLoad(b.CreateBitCast(f.FirstParam(), llir.PointerType(t), ""), ""))

	return f
}

func (g *moduleGenerator) environmentToEntryFunctionPointer(
	b llvm.Builder, v llvm.Value, t llvm.Type,
) llvm.Value {
	return b.CreateGEP(
		b.CreateBitCast(
			v,
			llir.PointerType(
				llir.PointerType(
					llir.FunctionType(t, []llvm.Type{types.NewEnvironment(0).LLVMPointerType()}),
				),
			),
			"",
		),
		[]llvm.Value{llvm.ConstIntFromString(llvm.Int32Type(), "-1", 10)},
		"",
	)
}

func (g moduleGenerator) lambdaToEnvironment(l ast.Lambda) types.Environment {
	if l.IsUpdatable() {
		return types.NewEnvironment(typeSize(g.module, types.Unbox(l.ResultType()).LLVMType()))
	}

	return types.NewEnvironment(0)
}

func (g moduleGenerator) createLogicalEnvironment(f llvm.Value, b llvm.Builder, l ast.Lambda) map[string]llvm.Value {
	vs := copyVariables(g.globalVariables)

	e := b.CreateBitCast(
		f.FirstParam(),
		llir.PointerType(lambdaToFreeVariablesStructType(l)),
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
