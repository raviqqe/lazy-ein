package types

// Float64 is a 64-bit float type.
type Float64 struct{}

// NewFloat64 creates a 64-bit float type.
func NewFloat64() Float64 {
	return Float64{}
}

// ConvertTypes converts types.
func (f Float64) ConvertTypes(ff func(Type) Type) Type {
	return ff(f)
}

func (Float64) String() string {
	return "Float64"
}

func (Float64) equal(t Type) bool {
	_, ok := t.(Float64)
	return ok
}
