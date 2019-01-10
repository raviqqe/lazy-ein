package types

import "fmt"

// Algebraic is an algebraic type.
type Algebraic struct {
	constructors []Constructor
}

// NewAlgebraic creates an algebraic type.
func NewAlgebraic(cs []Constructor) Algebraic {
	return Algebraic{cs}
}

// Constructors returns constructors.
func (a Algebraic) Constructors() []Constructor {
	return a.constructors
}

func (a Algebraic) String() string {
	s := a.constructors[0].String()

	for _, c := range a.constructors[1:] {
		s += "," + c.String()
	}

	return fmt.Sprintf("Algebraic([%v])", s)
}

func (Algebraic) isType() {}
