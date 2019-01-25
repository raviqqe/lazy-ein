package types

import "fmt"

// Index is an index referring to a recursive parent type like De Bruijn index.
// Its value starts from zero.
type Index struct {
	value int
}

// NewIndex creates a new index.
func NewIndex(i int) Index {
	return Index{i}
}

// Value returns a value.
func (i Index) Value() int {
	return i.value
}

// ConvertTypes converts types.
func (i Index) ConvertTypes(f func(Type) Type) Type {
	return f(i)
}

func (i Index) String() string {
	return fmt.Sprintf("&%v", i.value)
}

// equal checks structural equality.
func (i Index) equal(t Type) bool {
	ii, ok := t.(Index)

	if !ok {
		return false
	}

	return i.value == ii.value
}

func (Index) isBoxable() {}
