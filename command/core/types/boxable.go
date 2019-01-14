package types

// Boxable is a boxable type.
type Boxable interface {
	Type
	isBoxable()
}
