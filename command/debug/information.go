package debug

import "fmt"

// Information is debug information.
type Information struct {
	filename       string
	line, position int
	code           string
}

// NewInformation creates debug information.
func NewInformation(f string, l, p int, c string) *Information {
	return &Information{f, l, p, c}
}

func (i *Information) String() string {
	return fmt.Sprintf("%s:%d:%d:\t%s", i.filename, i.line, i.position, i.code)
}
