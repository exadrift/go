package tui

import "fmt"

type Dimensions struct {
	Top    int
	Left   int
	Width  int
	Height int
}

type Widget interface {
	GetChildren() []Widget
	GetDimensions() *Dimensions
	Render(mode RenderMode, inFocus Widget)
	SetDimensions(left int, top int, width int, height int)
	GetBox() *Box
	CaptureInput(r string) string
	FindInFocus(me Widget, prevInFocus Widget, inFocus Widget, focus *Focus) *Focus
	CanHaveFocus() bool
}

type Focus struct {
	First Widget
	Next  Widget
	Me    Widget
	Prev  Widget
	Last  Widget
}

type OptionType int

const (
	SegmentOptionMinChars OptionType = iota
	SegmentOptionMaxChars
	ApplicationOptionWithOnStart
	ApplicationOptionWithOnExit
	ApplicationOptionExitSignals
	ApplicationOptionInputHandler
)

type Option struct {
	optionType OptionType
	data       any
}

// Constrain will constrain the provided value to the provided length, adding ellipsis
// to the end if possible
func Constrain(value string, length int) string {
	if len(value) <= length {
		return value
	}

	newValue := []rune(value[:length])
	pos := len(newValue) - 1
	last := pos - 3
	if last < 0 {
		last = 0
	}
	for i := pos; i >= last; i-- {
		if pos >= 0 {
			newValue[pos] = '.'
		}
	}

	return string(newValue)
}

func Pad(value string, length int) string {
	return fmt.Sprintf("%-*s", length, value)
}
