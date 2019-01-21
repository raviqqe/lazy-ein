package types

import "fmt"

// Boxed is a boxed type.
type Boxed struct {
	algebraic Algebraic
}

// NewBoxed creates a boxed type.
func NewBoxed(t Algebraic) Boxed {
	return Boxed{t}
}

// Content returns a algebraic type.
func (b Boxed) Content() Algebraic {
	return b.algebraic
}

// ConvertTypes converts types.
func (b Boxed) ConvertTypes(f func(Type) Type) Type {
	return f(Boxed{b.algebraic.ConvertTypes(f).(Algebraic)})
}

func (b Boxed) String() string {
	return fmt.Sprintf("Boxed(%v)", b.algebraic)
}

func (b Boxed) equal(t Type) bool {
	bb, ok := t.(Boxed)

	if !ok {
		return false
	}

	return b.algebraic.equal(bb.algebraic)
}

func (Boxed) isBindable() {}
