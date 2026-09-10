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
	assert.Equal(t, 8, len(text.text))
	assert.Equal(t, []rune("hi there"), text.text[7])
}

func TestTextNewlineBeginning(t *testing.T) {
	origStr := "\nhello world\n\n1\nthis is a new text\nhi there"
	text := T(origStr)
	strippedText := strings.ReplaceAll(origStr, "\n", "")
	assert.Equal(t, text.Len(), len(strippedText))
	assert.Equal(t, 9, len(text.text))
	assert.Equal(t, Break, text.text[0])
}

func TestTextNewlineEnding(t *testing.T) {
	origStr := "hello world\n\n1\nthis is a new text\nhi there\n"
	text := T(origStr)
	strippedText := strings.ReplaceAll(origStr, "\n", "")
	assert.Equal(t, text.Len(), len(strippedText))
	assert.Equal(t, 9, len(text.text))
	assert.Equal(t, Break, text.text[8])
}

func TestTextRenderNoConstraints(t *testing.T) {
	text := T("hello there\nthis is a rendering")
	strs := text.Render()
	assert.Len(t, strs, 2)
}

func TestTextRenderMinumumOneRow(t *testing.T) {
	text := T().Render()
	assert.Len(t, text, 1)
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

func TestRenderWithStyleUnderrides(t *testing.T) {
	text := T("please render my text")
	s := text.Render(WithDefaultStyles(Blue.Fg()))
	assert.True(t, strings.HasPrefix(s[0], Blue.Fg().Ansi))

}

func TestRenderWithStyleReset(t *testing.T) {
	text := T("please render my", StyleReset, " text")
	s := text.Render(WithDefaultStyles(Blue.Fg()))
	assert.True(t, strings.HasPrefix(s[0], Blue.Fg().Ansi))
	assert.True(t, strings.HasSuffix(s[0], Blue.Fg().Ansi+" text"))
}

func TestExtendText(t *testing.T) {
	text := T("hello")
	text = text.Extend(T("world"))
	assert.Len(t, text.text, 2)
}

func TestExtendString(t *testing.T) {
	text := T("hello")
	text = text.Extend("world")
	assert.Len(t, text.text, 2)
}

func TestRequireScrollWidth(t *testing.T) {
	text := T("we have some text which will")
	reqScroll := text.RequiresScroll(5, 5)
	assert.True(t, reqScroll)

	reqScroll = text.RequiresScroll(5, 6)
	assert.False(t, reqScroll)
}

func TestRequireScrollWidthMultiPart(t *testing.T) {
	text := T("we have some", " text which will")
	reqScroll := text.RequiresScroll(5, 5)
	assert.True(t, reqScroll)

	reqScroll = text.RequiresScroll(5, 6)
	assert.False(t, reqScroll)
}

func TestRequireScrollWidthNewline(t *testing.T) {
	text := T("we have some\n", " text which will.")
	reqScroll := text.RequiresScroll(5, 5)
	assert.True(t, reqScroll)

	reqScroll = text.RequiresScroll(5, 6)
	assert.True(t, reqScroll)
}

func TestWrapText(t *testing.T) {
	text := T("some ", "text needs ", "to be wrapped")
	rows := text.Wrap(5)
	assert.Len(t, rows, 6)
}

func TestWrapTextWithNewline(t *testing.T) {
	text := T("some ", "text needs ", "to be wrapped\n\n")
	rows := text.Wrap(5)
	assert.Len(t, rows, 7)
}

func TestWrapTextWithMidNewline(t *testing.T) {
	text := T("some ", "text needs \n", "to be wrapped\n\n")
	rows := text.Wrap(5)
	assert.Len(t, rows, 8)
}
