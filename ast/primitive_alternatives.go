package ast

// PrimitiveAlternatives are primitive alternatives.
type PrimitiveAlternatives []PrimitiveAlternative

// NewPrimitiveAlternatives creates primitive alternatives.
func NewPrimitiveAlternatives(as ...PrimitiveAlternative) PrimitiveAlternatives {
	return as
}

func (as PrimitiveAlternatives) toSlice() []Alternative {
	aas := make([]Alternative, 0, len(as))

	for _, a := range as {
		aas = append(aas, Alternative(a))
	}

	return aas
}
