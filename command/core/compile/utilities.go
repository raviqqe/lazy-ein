package compile

import (
	"github.com/ein-lang/ein/command/core/compile/llir"
	"github.com/llvm-mirror/llvm/bindings/go/llvm"
)

func forceThunk(b llvm.Builder, thunk llvm.Value, g typeGenerator) llvm.Value {
	f := b.CreateBitCast(
		b.CreateCall(
			b.GetInsertBlock().Parent().GlobalParent().NamedFunction(atomicLoadFunctionName),
			[]llvm.Value{
				b.CreateBitCast(
					b.CreateStructGEP(thunk, 0, ""),
					llir.PointerType(llir.PointerType(llvm.Int8Type())),
					"",
				),
			},
			"",
		),
		thunk.Type().ElementType().StructElementTypes()[0],
		"",
	)

	wait := llvm.AddBasicBlock(b.GetInsertBlock().Parent(), "force.wait")
	lock := llvm.AddBasicBlock(b.GetInsertBlock().Parent(), "force.lock")
	update := llvm.AddBasicBlock(b.GetInsertBlock().Parent(), "force.update")
	end := llvm.AddBasicBlock(b.GetInsertBlock().Parent(), "force.end")

	b.CreateCondBr(
		b.CreateICmp(
			llvm.IntEQ,
			b.CreatePtrToInt(f, llir.WordType(), ""),
			b.CreatePtrToInt(llvm.ConstNull(f.Type()), llir.WordType(), ""),
			"",
		),
		wait,
		lock,
	)

	b.SetInsertPointAtEnd(wait)
	b.CreateCall(
		b.GetInsertBlock().Parent().GlobalParent().NamedFunction(blackHoleFunctionName),
		[]llvm.Value{b.CreateBitCast(thunk, llir.PointerType(llvm.Int8Type()), "")},
		"",
	)
	result1 := callEntryFunction(b, g, b.CreateLoad(b.CreateStructGEP(thunk, 0, ""), ""), thunk)
	b.CreateBr(end)

	b.SetInsertPointAtEnd(lock)
	b.CreateCondBr(
		b.CreateCall(
			b.GetInsertBlock().Parent().GlobalParent().NamedFunction(atomicCmpxchgFunctionName),
			[]llvm.Value{
				b.CreateBitCast(
					b.CreateStructGEP(thunk, 0, ""),
					llir.PointerType(llir.PointerType(llvm.Int8Type())),
					"",
				),
				b.CreateBitCast(f, llir.PointerType(llvm.Int8Type()), ""),
				llvm.ConstNull(llir.PointerType(llvm.Int8Type())),
			},
			"",
		),
		update,
		wait,
	)

	b.SetInsertPointAtEnd(update)
	result2 := callEntryFunction(b, g, f, thunk)
	// for cases that entry functions are not updated
	b.CreateCall(
		b.GetInsertBlock().Parent().GlobalParent().NamedFunction(atomicCmpxchgFunctionName),
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(thunk, 0, ""),
				llir.PointerType(llir.PointerType(llvm.Int8Type())),
				"",
			),
			llvm.ConstNull(llir.PointerType(llvm.Int8Type())),
			b.CreateBitCast(f, llir.PointerType(llvm.Int8Type()), ""),
		},
		"",
	)
	b.CreateBr(end)

	b.SetInsertPointAtEnd(end)
	p := b.CreatePHI(f.Type().ElementType().ReturnType(), "")
	p.AddIncoming([]llvm.Value{result1, result2}, []llvm.BasicBlock{wait, update})
	return p
}

func callEntryFunction(b llvm.Builder, g typeGenerator, f, thunk llvm.Value) llvm.Value {
	return b.CreateCall(
		f,
		[]llvm.Value{
			b.CreateBitCast(
				b.CreateStructGEP(thunk, 1, ""),
				llir.PointerType(g.GenerateUnsizedPayload()),
				"",
			),
		},
		"",
	)
}
