package types

// TypeVisitor visits types
type TypeVisitor interface {
	VisitTypes(func(Type) error) error
}
