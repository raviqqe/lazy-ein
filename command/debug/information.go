package debug

import "fmt"

// Information is debug information.
type Information struct {
	filename     string
	line, column int
	source       string
}

// NewInformation creates debug information.
func NewInformation(f string, l, c int, s string) *Information {
	return &Information{f, l, c, s}
}

func (i *Information) String() string {
	return fmt.Sprintf("%s:%d:%d:\t%s", i.filename, i.line, i.column, i.source)
}
