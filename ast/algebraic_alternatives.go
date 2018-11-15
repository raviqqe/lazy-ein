package ast

// AlgebraicAlternatives are algebraic alternatives.
type AlgebraicAlternatives []AlgebraicAlternative

// NewAlgebraicAlternatives creates algebraic alternatives.
func NewAlgebraicAlternatives(as ...AlgebraicAlternative) AlgebraicAlternatives {
	return as
}

func (as AlgebraicAlternatives) toSlice() []Alternative {
	aas := make([]Alternative, 0, len(as))

	for _, a := range as {
		aas = append(aas, Alternative(a))
	}

	return aas
}
