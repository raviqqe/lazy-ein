package ast

// PrimitiveAlternatives are primitive alternatives.
type PrimitiveAlternatives struct {
	alternatives []Alternative
}

// NewPrimitiveAlternatives creates primitive alternatives.
func NewPrimitiveAlternatives(as ...PrimitiveAlternative) PrimitiveAlternatives {
	aas := make([]Alternative, 0, len(as))

	for _, a := range as {
		aas = append(aas, Alternative(a))
	}

	return PrimitiveAlternatives{aas}
}

func (as PrimitiveAlternatives) toSlice() []Alternative {
	return as.alternatives
}
