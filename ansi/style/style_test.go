package style

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTextNewline(t *testing.T) {
	origStr := "hello world\n\n1\nthis is a new text\nhi there"
	text := T(origStr)
	strippedText := strings.ReplaceAll(origStr, "\n", "")
	assert.Equal(t, text.Len(), len(strippedText))
	assert.Equal(t, 8, len(text))
	assert.Equal(t, []rune("hi there"), text[7])
}

func TestTextNewlineBeginning(t *testing.T) {
	origStr := "\nhello world\n\n1\nthis is a new text\nhi there"
	text := T(origStr)
	strippedText := strings.ReplaceAll(origStr, "\n", "")
	assert.Equal(t, text.Len(), len(strippedText))
	assert.Equal(t, 9, len(text))
	assert.Equal(t, Break, text[0])
}

func TestTextNewlineEnding(t *testing.T) {
	origStr := "hello world\n\n1\nthis is a new text\nhi there\n"
	text := T(origStr)
	strippedText := strings.ReplaceAll(origStr, "\n", "")
	assert.Equal(t, text.Len(), len(strippedText))
	assert.Equal(t, 9, len(text))
	assert.Equal(t, Break, text[8])
}

func TestTextRenderNoConstraints(t *testing.T) {
	text := T("hello there\nthis is a rendering")
	strs := text.Render()
	assert.Len(t, strs, 2)
}

func TestTextRenderWithSomeAnsi(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n%s", line1, line2))
	strs := text.Render()
	assert.Len(t, strs, 2)
	assert.Len(t, StripAnsi(strs[0]), len(line1))
	assert.Len(t, StripAnsi(strs[1]), len(line2))
}

func TestTextRenderWithSomeAnsiAndDoubleNewline(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n\n%s", line1, line2))
	strs := text.Render()
	assert.Len(t, strs, 3)
	assert.Len(t, StripAnsi(strs[0]), len(line1))
	assert.Len(t, StripAnsi(strs[2]), len(line2))
}

func TestTextRenderWithSomeAnsiTrailingNewline(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n%s\n", line1, line2))
	strs := text.Render()
	assert.Len(t, strs, 2)
	assert.Len(t, StripAnsi(strs[0]), len(line1))
	assert.Len(t, StripAnsi(strs[1]), len(line2))
}

func TestTextRenderPadWidthWrap(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n%s\n", line1, line2))
	strs := text.Render(WithWidthConstraint(8))
	assert.Len(t, strs, 5)
	assert.Equal(t, StripAnsi(strs[0]), "hello th")
	assert.Equal(t, StripAnsi(strs[1]), "ere     ")
	assert.Equal(t, StripAnsi(strs[2]), "this is ")
	assert.Equal(t, StripAnsi(strs[3]), "a render")
	assert.Equal(t, StripAnsi(strs[4]), "ing     ")
}

func TestTextRenderPadWidthWrapMinRows(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n%s\n", line1, line2))
	strs := text.Render(WithWidthConstraint(8), WithMinRows(10))
	assert.Len(t, strs, 10)
	assert.Equal(t, StripAnsi(strs[0]), "hello th")
	assert.Equal(t, StripAnsi(strs[1]), "ere     ")
	assert.Equal(t, StripAnsi(strs[2]), "this is ")
	assert.Equal(t, StripAnsi(strs[3]), "a render")
	assert.Equal(t, StripAnsi(strs[4]), "ing     ")

	for _, row := range strs {
		assert.Len(t, StripAnsi(row), 8)
	}
}

func TestTextRenderPadExtra(t *testing.T) {
	line1 := "hello there"
	line2 := "this is a rendering"
	text := T(Blue.Fg(), fmt.Sprintf("%s\n%s\n", line1, line2))
	strs := text.Render(WithWidthConstraint(20))
	assert.Len(t, strs, 2)
	assert.Equal(t, StripAnsi(strs[0]), "hello there         ")
	assert.Equal(t, StripAnsi(strs[1]), "this is a rendering ")
}

func TestApplyDefaultStyles(t *testing.T) {
	text := T("hello world")
	restyled := text.WithDefaultStyles(Blue.Fg())
	assert.Len(t, restyled, 2)
}

func TestApplyDefaultStylesTwo(t *testing.T) {
	text := T(Red.Bg(), "hello world")
	restyled := text.WithDefaultStyles(Blue.Fg())
	assert.Len(t, restyled, 3)
	assert.Equal(t, Blue.Fg(), restyled[0].(Style))
}

func TestApplyDefaultStylesUnderride(t *testing.T) {
	text := T(Red.Fg(), "hello world", StyleReset, "normal text")
	restyled := text.WithDefaultStyles(Blue.Fg())
	assert.Len(t, restyled, 6)
	assert.Equal(t, Blue.Fg(), restyled[0].(Style))

	// this gets inserted right before "normal text"
	assert.Equal(t, Blue.Fg(), restyled[4].(Style))
}

func TestExtendText(t *testing.T) {
	text := T("hello")
	text = text.Extend(T("world"))
	assert.Len(t, text, 2)
}

func TestExtendString(t *testing.T) {
	text := T("hello")
	text = text.Extend("world")
	assert.Len(t, text, 2)
}
