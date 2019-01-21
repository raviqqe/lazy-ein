package types

import "fmt"

// Boxed is a boxed type.
type Boxed struct {
	content Boxable
}

// NewBoxed creates a boxed type.
func NewBoxed(t Boxable) Boxed {
	return Boxed{t}
}

// Content returns a content type.
func (b Boxed) Content() Boxable {
	return b.content
}

// ConvertTypes converts types.
func (b Boxed) ConvertTypes(f func(Type) Type) Type {
	return f(Boxed{b.content.ConvertTypes(f).(Boxable)})
}

func (b Boxed) String() string {
	return fmt.Sprintf("Boxed(%v)", b.content)
}

func (b Boxed) equal(t Type) bool {
	bb, ok := t.(Boxed)

	if !ok {
		return false
	}

	return b.content.equal(bb.content)
}

func (Boxed) isBindable() {}
