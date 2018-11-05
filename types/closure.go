package types

import (
	"llvm.org/llvm/bindings/go/llvm"
)

// Closure is a closure type.
type Closure struct {
	environment Environment
	arguments   []Type
	result      Type
}

// NewClosure creates a new closure type.
func NewClosure(e Environment, as []Type, r Type) Closure {
	return Closure{e, as, r}
}

// LLVMType returns a LLVM type.
func (c Closure) LLVMType() llvm.Type {
	return llvm.StructType(
		[]llvm.Type{
			llvm.PointerType(
				llvm.FunctionType(
					c.result.LLVMType(),
					append(
						[]llvm.Type{NewEnvironment(0).LLVMPointerType()},
						ToLLVMTypes(c.arguments)...,
					),
					false,
				),
				0,
			),
			c.environment.LLVMType(),
		},
		false,
	)
}
