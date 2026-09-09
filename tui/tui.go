package tui

import (
	"fmt"

	"github.com/exadrift/go/ansi/style"
)

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
	Collect(me Widget) []Widget
	SetDimensions(left int, top int, width int, height int)
	GetBox() *Box
	CaptureInput(r string) string
	GetFocalWidgets(me Widget, focalWidgets *FocalWidgets)
	CanHaveFocus() bool
	AbsorbsInput(input string) bool
	ResetScrollPosition()
}

type FocalWidgets struct {
	Widgets []Widget
}

type OptionType int

const (
	SegmentOptionMinChars OptionType = iota
	SegmentOptionMaxChars
	ApplicationOptionWithOnStart
	ApplicationOptionWithOnExit
	ApplicationOptionExitSignals
	ApplicationOptionInputHandler
	ApplicationOptionKeyBindings
	BusyModal
)

type Option struct {
	optionType OptionType
	data       any
}

type KeyBindings struct {
	FocusNext     string
	FocusPrev     string
	SelectionNext string
	SelectionPrev string
	ScrollUp      string
	ScrollDown    string
	Trigger       string
}

func NewKeyBindings() *KeyBindings {
	return &KeyBindings{
		FocusNext:     Tab,
		FocusPrev:     ShiftTab,
		SelectionNext: DownArrow,
		SelectionPrev: UpArrow,
		ScrollUp:      CtrlPgUp,
		ScrollDown:    CtrlPgDn,
		Trigger:       Enter,
	}
}

// Constrain will constrain the provided value to the provided length, adding ellipsis
// to the end if possible
func Constrain(value string, length int) string {
	if len(value) <= length {
		return value
	}

	newValue := []rune(value[:length])
	last := len(newValue) - 1
	first := last - 1
	if first < 0 {
		first = 0
	}
	for i := first; i <= last; i++ {
		newValue[i] = '.'
	}

	return string(newValue)
}

func Pad(value string, length int) string {
	return fmt.Sprintf("%-*s", length, value)
}

func ProcessToStyledText(text any) style.Text {
	switch t := text.(type) {
	case style.Text:
		return t
	case []rune:
		return style.T(string(t))
	case string:
		return style.T(t)
	default:
		panic(fmt.Sprintf("unknown text type to be processed into stylized text %+v", text))
	}
}
