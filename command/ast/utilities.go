package ast

// PatternsEqual checks equality of patterns.
func PatternsEqual(e, ee Expression) bool {
	switch e := e.(type) {
	case List:
		l, ok := ee.(List)

		if !ok || len(e.Arguments()) != len(l.Arguments()) {
			return false
		} else if len(e.Arguments()) == 0 {
			return true
		} else if e.Arguments()[0].Expanded() != l.Arguments()[0].Expanded() {
			return false
		}

		return PatternsEqual(e.Arguments()[0].Expression(), l.Arguments()[0].Expression()) &&
			PatternsEqual(NewList(e.Type(), e.Arguments()[1:]), NewList(l.Type(), l.Arguments()[1:]))
	case Number:
		n, ok := ee.(Number)
		return ok && e.Value() == n.Value()
	case Variable:
		_, ok := ee.(Variable)
		return ok
	}

	panic("unreachable")
}
