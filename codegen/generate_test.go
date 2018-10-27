package codegen

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

var environmentOffset = reflect.PtrTo(reflect.TypeOf(42)).Size()

func TestGenerate(t *testing.T) {
	m, err := Generate(
		"foo",
		[]ast.Bind{
			ast.NewBind("foo", ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
		},
	)

	assert.NotEqual(t, llvm.Module{}, m)
	assert.Nil(t, err)
}

func TestGlobalThunkForce(t *testing.T) {
	const functionName = "foo"

	m, err := Generate(
		"foo",
		[]ast.Bind{
			ast.NewBind(functionName, ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64())),
		},
	)
	assert.Nil(t, err)

	e, err := llvm.NewExecutionEngine(m)
	assert.Nil(t, err)

	g := e.PointerToGlobal(m.NamedGlobal(functionName))
	p := unsafe.Pointer(uintptr(g) + environmentOffset)

	assert.NotEqual(t, 42.0, *(*float64)(unsafe.Pointer(uintptr(g) + environmentOffset)))

	assert.Equal(t, 42.0, e.RunFunction(
		e.FindFunction(toEntryName(functionName)),
		[]llvm.GenericValue{
			llvm.NewGenericValueFromPointer(p),
		},
	).Float(llvm.DoubleType()))

	assert.Equal(t, 42.0, *(*float64)(unsafe.Pointer(uintptr(g) + environmentOffset)))

	assert.Equal(t, 42.0, e.RunFunction(
		e.FindFunction(toUpdatedEntryName(functionName)),
		[]llvm.GenericValue{
			llvm.NewGenericValueFromPointer(p),
		},
	).Float(llvm.DoubleType()))
}