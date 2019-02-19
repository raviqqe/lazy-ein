package ast

// Variable is a variable.
type Variable struct {
	name string
}

// NewVariable creates a variable.
func NewVariable(s string) Variable {
	return Variable{s}
}

// Name returns a name.
func (v Variable) Name() string {
	return v.name
}

// RenameVariablesInAtom renames variables.
func (v Variable) RenameVariablesInAtom(vs map[string]string) Atom {
	return Variable{replaceVariable(v.name, vs)}
}

func (Variable) isAtom() {}
