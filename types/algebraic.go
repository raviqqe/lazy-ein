package types

// Algebraic is an algebraic type.
type Algebraic struct {
	constructors []Constructor
}

// NewAlgebraic creates a new algebraic type.
func NewAlgebraic(cs []Constructor) Algebraic {
	return Algebraic{cs}
}

// Constructors returns constructors.
func (a Algebraic) Constructors() []Constructor {
	return a.constructors
}

func (Algebraic) isType() {}
