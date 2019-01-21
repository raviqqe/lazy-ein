package types

// Bindable is a bindable type.
type Bindable interface {
	Type
	isBindable()
}
