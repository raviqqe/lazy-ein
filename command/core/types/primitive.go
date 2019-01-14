package types

// Primitive is a primitive type.
type Primitive interface {
	Type
	isPrimitive()
}
