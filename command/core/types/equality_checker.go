package types

type equalityChecker struct {
	pairs                 [][2]Type
	leftStack, rightStack []Type
}

func newEqualityChecker() equalityChecker {
	return equalityChecker{nil, nil, nil}
}

func (c equalityChecker) Check(t, tt Type) bool {
	if c.isPairChecked(t, tt) {
		return true
	}

	c = c.addPair(t, tt)

	if i, ok := t.(Index); ok {
		return c.Check(c.leftStack[len(c.leftStack)-1-i.Value()], tt)
	} else if i, ok := tt.(Index); ok {
		return c.Check(t, c.rightStack[len(c.rightStack)-1-i.Value()])
	}

	switch t := t.(type) {
	case Algebraic:
		a, ok := tt.(Algebraic)

		if !ok || len(t.Constructors()) != len(a.Constructors()) {
			return false
		}

		c = c.pushTypes(t, tt)

		for i, cc := range t.Constructors() {
			es := cc.Elements()
			ees := a.Constructors()[i].Elements()

			if len(es) != len(ees) {
				return false
			}

			for i, e := range es {
				if !c.Check(e, ees[i]) {
					return false
				}
			}
		}

		return true
	}

	return t.equal(tt)
}

func (c equalityChecker) pushTypes(t, tt Type) equalityChecker {
	return equalityChecker{c.pairs, append(c.leftStack, t), append(c.rightStack, tt)}
}

func (c equalityChecker) addPair(t, tt Type) equalityChecker {
	return equalityChecker{append(c.pairs, [2]Type{t, tt}), c.leftStack, c.rightStack}
}

func (c equalityChecker) isPairChecked(t, tt Type) bool {
	for _, ts := range c.pairs {
		if t.equal(ts[0]) && tt.equal(ts[1]) {
			return true
		}
	}

	return false
}
