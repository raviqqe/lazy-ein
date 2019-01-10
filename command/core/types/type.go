package types

import "fmt"

// Type is a type.
type Type interface {
	fmt.Stringer
	isType()
}
