package types

import "fmt"

// Type is a type.
type Type interface {
	fmt.Stringer
	ConvertTypes(func(Type) Type) Type
	equal(Type) bool
}
