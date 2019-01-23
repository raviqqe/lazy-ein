package types

type recursivityChecker struct {
	stack []Type
}

func newRecursivityChecker() recursivityChecker {
	return recursivityChecker{nil}
}

func (c recursivityChecker) Check(t Type) bool {
	switch t := t.(type) {
	case Algebraic:
		c = c.pushType(t)

		for _, cc := range t.Constructors() {
			for _, e := range cc.Elements() {
				if c.Check(e) {
					return true
				}
			}
		}

		return false
	case Boxed:
		return c.Check(t.Content())
	case Function:
		c = c.pushType(t)

		for _, a := range t.Arguments() {
			if c.Check(a) {
				return true
			}
		}

		return c.Check(t.Result())
	case Index:
		return t.Value() == len(c.stack)-1
	}

	return false
}

func (c recursivityChecker) pushType(t Type) recursivityChecker {
	return recursivityChecker{append(c.stack, t)}
}
