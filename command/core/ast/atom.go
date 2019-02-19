package ast

// Atom is an atom.
type Atom interface {
	isAtom()
	RenameVariablesInAtom(map[string]string) Atom
}
