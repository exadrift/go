package style

import (
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ansiEscape = regexp.MustCompile(`\x1b(?:[@-Z\\-_]|\[[0-?]*[ -/]*[@-~])`)

func StripAnsi(value string) string {
	return ansiEscape.ReplaceAllString(value, "")
}

func TestTextLength(t *testing.T) {
	s1 := "hello world"
	s2 := " what's good?"
	text := T(S(s1, Blue.Fg()), S(s2, Red.Bg()))
	assert.Equal(t, len(s1)+len(s2), text.Len())
}

func TestTextRender(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	text := T(S(t1, s1), S(t2, s2))
	rendered, leftover := text.Render()
	assert.Nil(t, leftover)
	rendExp := s1.Ansi() + t1 + ResetStyleAnsi + s2.Ansi() + t2 + ResetStyleAnsi
	assert.Equal(t, rendExp, rendered)
}

func TestTextRenderStyleDefaults(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	text := T(S(t1, s1), S(t2, s2))
	rendered, leftover := text.Render(WithStyles(ds1, ds2))
	assert.Nil(t, leftover)
	rendExp := s1.Ansi() + ds2.Ansi() + t1 + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + t2 + ResetStyleAnsi
	assert.Equal(t, rendExp, rendered)
}

func TestTextRenderWithFixedWidthUnder(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 5

	text := T(S(t1, s1), S(t2, s2))
	rendered, leftover := text.Render(WithStyles(ds1, ds2), WithFixedWidth(width))
	assert.NotNil(t, leftover)
	rendExp := s1.Ansi() + ds2.Ansi() + t1[:width] + ResetStyleAnsi
	assert.Equal(t, rendExp, rendered)
	leftoverRend, leftover := leftover.Render(WithStyles(ds1, ds2))
	assert.Nil(t, leftover)
	leftOverRedExp := s1.Ansi() + ds2.Ansi() + t1[width:] + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + t2 + ResetStyleAnsi
	assert.Equal(t, leftOverRedExp, leftoverRend)
}

func TestTextRenderWithFixedWidthOver(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 100

	text := T(S(t1, s1), S(t2, s2))
	rendered, leftover := text.Render(WithStyles(ds1, ds2), WithFixedWidth(width))
	assert.Nil(t, leftover)
	rendExp := s1.Ansi() + ds2.Ansi() + t1 + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + t2 + ResetStyleAnsi + ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width-len(t1)-len(t2)) + ResetStyleAnsi
	assert.Equal(t, rendExp, rendered)
}

func TestTextBlock(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 10
	height := 5

	text := T(S(t1, s1), S(t2, s2))
	tb := B(text)
	lines := tb.Render(width, height, 0, ds1, ds2)
	assert.Len(t, lines, height)

	line1 := lines[0]
	expLine1 := s1.Ansi() + ds2.Ansi() + "hello worl" + ResetStyleAnsi
	assert.Equal(t, expLine1, line1)
	line2 := lines[1]
	expLine2 := s1.Ansi() + ds2.Ansi() + "d" + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + " what's g" + ResetStyleAnsi
	assert.Equal(t, expLine2, line2)
	line3 := lines[2]
	expLine3 := s2.Ansi() + ds1.Ansi() + "ood?" + ResetStyleAnsi + ds1.Ansi() + ds2.Ansi() + "      " + ResetStyleAnsi
	assert.Equal(t, expLine3, line3)
	line4 := lines[3]
	expLine4 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine4, line4)
	line5 := lines[4]
	expLine5 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine5, line5)
}

func TestTextBlockHeightConstrained(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 10
	height := 2

	text := T(S(t1, s1), S(t2, s2))
	tb := B(text)
	lines := tb.Render(width, height, 0, ds1, ds2)
	assert.Len(t, lines, height)

	line1 := lines[0]
	expLine1 := s1.Ansi() + ds2.Ansi() + "hello worl" + ResetStyleAnsi
	assert.Equal(t, expLine1, line1)
	line2 := lines[1]
	expLine2 := s1.Ansi() + ds2.Ansi() + "d" + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + " what's g" + ResetStyleAnsi
	assert.Equal(t, expLine2, line2)
}

func TestTextBlockHeightConstrainedPositiveOffset(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 10
	height := 5

	text := T(S(t1, s1), S(t2, s2))
	tb := B(text)
	lines := tb.Render(width, height, 1, ds1, ds2)
	assert.Len(t, lines, height)

	line2 := lines[0]
	expLine2 := s1.Ansi() + ds2.Ansi() + "d" + ResetStyleAnsi + s2.Ansi() + ds1.Ansi() + " what's g" + ResetStyleAnsi
	assert.Equal(t, expLine2, line2)
	line3 := lines[1]
	expLine3 := s2.Ansi() + ds1.Ansi() + "ood?" + ResetStyleAnsi + ds1.Ansi() + ds2.Ansi() + "      " + ResetStyleAnsi
	assert.Equal(t, expLine3, line3)
}

func TestTextBlockPositiveOffset(t *testing.T) {
	t1 := "hello world"
	t2 := " what's good?"
	s1 := Blue.Fg()
	s2 := Red.Bg()

	ds1 := Green.Fg()
	ds2 := Green.Bg()

	width := 10
	height := 5

	text := T(S(t1, s1), S(t2, s2))
	tb := B(text)
	lines := tb.Render(width, height, 2, ds1, ds2)
	assert.Len(t, lines, height)

	line3 := lines[0]
	expLine3 := s2.Ansi() + ds1.Ansi() + "ood?" + ResetStyleAnsi + ds1.Ansi() + ds2.Ansi() + "      " + ResetStyleAnsi
	assert.Equal(t, expLine3, line3)
	line4 := lines[1]
	expLine4 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine4, line4)
	line5 := lines[2]
	expLine5 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine5, line5)
	line6 := lines[3]
	expLine6 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine6, line6)
	line7 := lines[4]
	expLine7 := ds1.Ansi() + ds2.Ansi() + strings.Repeat(" ", width) + ResetStyleAnsi
	assert.Equal(t, expLine7, line7)
}

func TestStyleTypeRemove(t *testing.T) {
	styles := Styles{Black.Bg(), Green.Fg()}
	styles = styles.Remove(StyleTypeBgColor)
	assert.Len(t, styles, 1)
	assert.Equal(t, styles[0].styleType, StyleTypeFgColor)
}

func TestStyleTypeRemoveBoth(t *testing.T) {
	styles := Styles{Black.Bg(), Green.Fg()}
	styles = styles.Remove(StyleTypeBgColor, StyleTypeFgColor)
	assert.Len(t, styles, 0)
}

func TestStyleTypeAdd(t *testing.T) {
	styles := Styles{Black.Bg(), Green.Fg()}
	styles = styles.Add(Red.Fg())
	assert.Len(t, styles, 2)
	assert.Equal(t, styles[0].styleType, StyleTypeBgColor)
	assert.Equal(t, styles[1].styleType, StyleTypeFgColor)
	assert.Equal(t, styles[1].ansi, Red.Fg().Ansi())
}

func TestCountLines(t *testing.T) {
	width := 10

	s1 := "hello world, how are you"
	s2 := "we need to render these lines"
	t1 := T(s1)
	t2 := T(s2)
	tb := B(t1, t2)
	numLines := tb.NumLines(width)
	numLinesT1 := len(s1) / width
	if len(s1)%width > 0 {
		numLinesT1++
	}
	numLinesT2 := len(s2) / width
	if len(s2)%width > 0 {
		numLinesT2++
	}

	assert.Equal(t, numLinesT1+numLinesT2, numLines)
}

func TestStyleBasicCount(t *testing.T) {
	s := S("hello", Blue.Fg(), Blue.Bg())
	assert.Len(t, s.styles, 2)
}

func TestStyleRgnCount(t *testing.T) {
	s := S("hello", FromRgb(10, 10, 10).Fg(), FromRgb(10, 10, 10).Bg())
	assert.Len(t, s.styles, 2)
}

func TestEmptyTextBlock(t *testing.T) {
	b := B()
	width := 200
	height := 1
	rend := b.Render(width, height, 0)
	assert.Len(t, rend, height)
	for _, line := range rend {
		assert.Len(t, StripAnsi(line), width)
	}
}
