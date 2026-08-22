package tui

import (
	"fmt"
	"sync"

	"github.com/exadrift/go/tui/internal/terminal"
)

type Loader struct {
	*Box
	Label string
	lock  sync.Mutex
}

func NewLoader(label string) *Loader {
	return &Loader{
		Box:   NewBox(),
		Label: label,
	}
}

func (l *Loader) SetLabel(label string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Label = label
}

func (l *Loader) Render(mode RenderMode, focusItem Widget) {
	l.lock.Lock()
	defer l.lock.Unlock()

	width := l.GetWidth()

	dimensions := l.GetDimensions()
	curY := dimensions.Top

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	for i := 0; i < width; i++ {
		fmt.Print(" ")
	}
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	fmt.Print(" ")
	fmt.Print(StyleFg(Blue))
	for i := 1; i < width-1; i++ {
		fmt.Print("█")
	}
	fmt.Print(StyleReset)
	fmt.Print(" ")
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	fmt.Print(" ")
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleReset)
	for i := 2; i < width-2; i++ {
		fmt.Print(" ")
	}
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleReset)
	fmt.Print(" ")
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	fmt.Print(" ")
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleFg(Yellow))
	fmt.Printf(" %s ", l.Label)
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleReset)
	fmt.Print(" ")
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	fmt.Print(" ")
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleReset)
	for i := 2; i < width-2; i++ {
		fmt.Print(" ")
	}
	fmt.Print(StyleFg(Blue))
	fmt.Print("█")
	fmt.Print(StyleReset)
	fmt.Print(" ")
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	fmt.Print(" ")
	fmt.Print(StyleFg(Blue))
	for i := 1; i < width-1; i++ {
		fmt.Print("█")
	}
	fmt.Print(StyleReset)
	fmt.Print(" ")
	curY++

	terminal.SetCursorPos(dimensions.Left, curY)
	fmt.Print(StyleReset)
	for i := 0; i < width; i++ {
		fmt.Print(" ")
	}
	curY++
}

func (l *Loader) GetWidth() int {
	return len(l.Label) + 4
}

func (l *Loader) GetHeight() int {
	return 7
}
