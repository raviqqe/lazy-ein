package types

// Algebraic is an algebraic type.
type Algebraic struct {
	name         string
	constructors []Constructor
}

// NewAlgebraic creates an algebraic type.
func NewAlgebraic(cs []Constructor) Algebraic {
	return Algebraic{"", cs}
}

// NewNamedAlgebraic creates an algebraic type.
func NewNamedAlgebraic(n string, cs []Constructor) Algebraic {
	return Algebraic{n, cs}
}

// Name returns a name.
func (a Algebraic) Name() (string, bool) {
	if a.name == "" {
		return "", false
	}

	return a.name, true
}

// Constructors returns constructors.
func (a Algebraic) Constructors() []Constructor {
	return a.constructors
}

func (Algebraic) isType() {}
