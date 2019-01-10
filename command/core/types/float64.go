package types

// Float64 is a 64-bit float type.
type Float64 struct{}

// NewFloat64 creates a 64-bit float type.
func NewFloat64() Float64 {
	return Float64{}
}

func (Float64) String() string {
	return "Float64"
}

func (Float64) isType() {}
