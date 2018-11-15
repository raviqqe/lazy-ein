package ast

// AlgebraicAlternatives are algebraic alternatives.
type AlgebraicAlternatives struct {
	alternatives []Alternative
}

// NewAlgebraicAlternatives creates algebraic alternatives.
func NewAlgebraicAlternatives(as ...AlgebraicAlternative) AlgebraicAlternatives {
	aas := make([]Alternative, 0, len(as))

	for _, a := range as {
		aas = append(aas, Alternative(a))
	}

	return AlgebraicAlternatives{aas}
}

func (as AlgebraicAlternatives) toSlice() []Alternative {
	return as.alternatives
}
