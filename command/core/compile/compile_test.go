package compile_test

import (
	"reflect"
	"testing"
	"unsafe"

	"github.com/ein-lang/ein/command/core/ast"
	"github.com/ein-lang/ein/command/core/compile"
	"github.com/ein-lang/ein/command/core/compile/names"
	"github.com/ein-lang/ein/command/core/types"
	"github.com/stretchr/testify/assert"
	"llvm.org/llvm/bindings/go/llvm"
)

var payloadOffset = reflect.PtrTo(reflect.TypeOf(42)).Size()
var algebraicType = types.NewAlgebraic(types.NewConstructor())

func TestCompile(t *testing.T) {
	m, err := compile.Compile(
		ast.NewModule(
			nil,
			[]ast.Bind{
				ast.NewBind(
					"foo",
					ast.NewVariableLambda(
						nil,
						true,
						ast.NewConstructorApplication(ast.NewConstructor(algebraicType, 0), nil),
						algebraicType,
					),
				),
			},
		),
	)

	assert.NotEqual(t, llvm.Module{}, m)
	assert.Nil(t, err)
}

func TestGlobalThunkForce(t *testing.T) {
	a := types.NewAlgebraic(types.NewConstructor(types.NewFloat64()))

	m, err := compile.Compile(
		ast.NewModule(
			nil,
			[]ast.Bind{
				ast.NewBind(
					"x",
					ast.NewVariableLambda(
						nil,
						true,
						ast.NewConstructorApplication(
							ast.NewConstructor(a, 0),
							[]ast.Atom{ast.NewFloat64(42)},
						),
						a,
					),
				),
			},
		),
	)
	assert.Nil(t, err)

	e, err := llvm.NewExecutionEngine(m)
	assert.Nil(t, err)

	p := unsafe.Pointer(uintptr(e.PointerToGlobal(m.NamedGlobal("x"))) + payloadOffset)

	assert.NotEqual(t, 42.0, *(*float64)(p))

	// TODO: Check return values in some way.
	e.RunFunction(
		e.FindFunction(names.ToEntry("x")),
		[]llvm.GenericValue{llvm.NewGenericValueFromPointer(p)},
	)

	assert.Equal(t, 42.0, *(*float64)(p))

	e.RunFunction(
		e.FindFunction(names.ToNormalFormEntry("x")),
		[]llvm.GenericValue{llvm.NewGenericValueFromPointer(p)},
	)

	assert.Equal(t, 42.0, *(*float64)(p))
}
