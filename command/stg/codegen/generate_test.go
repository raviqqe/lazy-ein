package codegen

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/raviqqe/stg/ast"
	"github.com/raviqqe/stg/codegen/names"
	"github.com/raviqqe/stg/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

var payloadOffset = reflect.PtrTo(reflect.TypeOf(42)).Size()

func TestGenerate(t *testing.T) {
	m, err := Generate(
		ast.NewModule(
			"foo",
			nil,
			[]ast.Bind{
				ast.NewBind(
					"foo",
					ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64()),
				),
			},
		),
	)

	assert.NotEqual(t, llvm.Module{}, m)
	assert.Nil(t, err)
}

func TestGlobalThunkForce(t *testing.T) {
	const functionName = "foo"

	m, err := Generate(
		ast.NewModule(
			"foo",
			nil,
			[]ast.Bind{
				ast.NewBind(
					functionName,
					ast.NewLambda(nil, true, nil, ast.NewFloat64(42), types.NewFloat64()),
				),
			},
		),
	)
	assert.Nil(t, err)

	e, err := llvm.NewExecutionEngine(m)
	assert.Nil(t, err)

	g := e.PointerToGlobal(m.NamedGlobal(functionName))
	p := unsafe.Pointer(uintptr(g) + payloadOffset)

	assert.NotEqual(t, 42.0, *(*float64)(unsafe.Pointer(uintptr(g) + payloadOffset)))

	assert.Equal(t, 42.0, e.RunFunction(
		e.FindFunction(names.ToEntry(functionName)),
		[]llvm.GenericValue{
			llvm.NewGenericValueFromPointer(p),
		},
	).Float(llvm.DoubleType()))

	assert.Equal(t, 42.0, *(*float64)(unsafe.Pointer(uintptr(g) + payloadOffset)))

	assert.Equal(t, 42.0, e.RunFunction(
		e.FindFunction(names.ToUpdatedEntry(functionName)),
		[]llvm.GenericValue{
			llvm.NewGenericValueFromPointer(p),
		},
	).Float(llvm.DoubleType()))
}
