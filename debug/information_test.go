package debug

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInformationString(t *testing.T) {
	assert.Equal(t, "filename:1:2:\tcode", fmt.Sprint(NewInformation("filename", 1, 2, "code")))
}
