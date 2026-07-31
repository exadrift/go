package tui

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstrainAnsiFullWidth(t *testing.T) {
	text := fmt.Sprintf("%sthis is some%s %stext with a bit%s of color as well", StyleFg(Red), StyleReset, StyleFgBg(Blue, White), StyleReset)
	constrained := ConstrainAnsiFullWidth(text, 10)
	stripped := StripAnsiCodes(constrained)
	assert.Len(t, stripped, 10)

	constrained = ConstrainAnsiFullWidth(text, 100)
	stripped = StripAnsiCodes(constrained)
	assert.Len(t, stripped, 100)
}
