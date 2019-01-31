package ast

// PatternsEqual checks equality of patterns.
func PatternsEqual(e, ee Expression) bool {
	switch e := e.(type) {
	case List:
		l, ok := ee.(List)

		if !ok || len(e.Elements()) != len(l.Elements()) {
			return false
		} else if len(e.Elements()) == 0 {
			return true
		}

		return PatternsEqual(e.Elements()[0], l.Elements()[0]) &&
			PatternsEqual(NewList(e.Type(), e.Elements()[1:]), NewList(l.Type(), l.Elements()[1:]))
	case Number:
		n, ok := ee.(Number)
		return ok && e.Value() == n.Value()
	}

	panic("unreachable")
}
