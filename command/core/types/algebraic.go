package types

import "fmt"

// Algebraic is an algebraic type.
type Algebraic struct {
	constructors []Constructor
}

// NewAlgebraic creates an algebraic type.
func NewAlgebraic(c Constructor, cs ...Constructor) Algebraic {
	return Algebraic{append([]Constructor{c}, cs...)}
}

// Constructors returns constructors.
func (a Algebraic) Constructors() []Constructor {
	return a.constructors
}

// ConvertTypes converts types.
func (a Algebraic) ConvertTypes(f func(Type) Type) Type {
	cs := make([]Constructor, 0, len(a.constructors))

	for _, c := range a.constructors {
		cs = append(cs, c.ConvertTypes(f))
	}

	return f(Algebraic{cs})
}

func (a Algebraic) String() string {
	s := a.constructors[0].String()

	for _, c := range a.constructors[1:] {
		s += "," + c.String()
	}

	return fmt.Sprintf("Algebraic(%v)", s)
}

func (a Algebraic) equal(t Type) bool {
	aa, ok := t.(Algebraic)

	if !ok {
		return false
	} else if len(a.constructors) != len(aa.constructors) {
		return false
	}

	for i, c := range a.constructors {
		if !c.equal(aa.constructors[i]) {
			return false
		}
	}

	return true
}

func (Algebraic) isBoxable() {}
