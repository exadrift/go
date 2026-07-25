package tui

import (
	"fmt"

	"github.com/exadrift/go/tui/internal/terminal"
)

type Alignment int

const (
	AlignmentLeft Alignment = iota
	AlignmentRight
	AlignmentCenter
)

const (
	BorderSingleTopLeftCorner     = "┌"
	BorderSingleHoriz             = "─"
	BorderSingleTopRightCorner    = "┐"
	BorderSingleVert              = "│"
	BorderSingleBottomLeftCorner  = "└"
	BorderSingleBottomRightCorner = "┘"

	BorderDoubleTopLeftCorner     = "╔"
	BorderDoubleHoriz             = "═"
	BorderDoubleTopRightCorner    = "╗"
	BorderDoubleVert              = "║"
	BorderDoubleBottomLeftCorner  = "╚"
	BorderDoubleBottomRightCorner = "╝"
)

type RenderMode int

const (
	RenderModeAll RenderMode = iota
	RenderModeBorder
	RenderModeContent
)

type Box struct {
	hasBorder         bool
	title             string
	alignment         Alignment
	dimensions        Dimensions
	contentDimensions Dimensions
}

func NewBox() *Box {
	return &Box{}
}

func (b *Box) Box() *Box {
	return b
}

func (b *Box) GetDimensions() *Dimensions {
	return &b.dimensions
}

func (b *Box) GetContentDimensions() *Dimensions {
	return &b.contentDimensions
}

func (b *Box) EnableBorder(enable bool) *Box {
	b.hasBorder = true
	return b
}

func (b *Box) SetDimensions(left int, top int, width int, height int) {
	b.dimensions.Left = left
	b.dimensions.Top = top
	b.dimensions.Width = width
	b.dimensions.Height = height

	if b.hasBorder {
		b.contentDimensions.Left = left + 1
		b.contentDimensions.Top = top + 1
		b.contentDimensions.Width = width - 2
		b.contentDimensions.Height = height - 2
	} else {
		b.contentDimensions.Left = left
		b.contentDimensions.Top = top
		b.contentDimensions.Width = width
		b.contentDimensions.Height = height
	}
}

func (b *Box) GetBox() *Box {
	return b
}

func (b *Box) Render(mode RenderMode, focusItem Widget) {
	if mode == RenderModeAll || mode == RenderModeBorder {
		if b.hasBorder && b.dimensions.Width >= 3 && b.dimensions.Height >= 3 {
			var (
				topLeft     string
				horiz       string
				topRight    string
				vert        string
				bottomLeft  string
				bottomRight string
			)

			if b == focusItem.GetBox() {
				topLeft = BorderDoubleTopLeftCorner
				horiz = BorderDoubleHoriz
				topRight = BorderDoubleTopRightCorner
				vert = BorderDoubleVert
				bottomLeft = BorderDoubleBottomLeftCorner
				bottomRight = BorderDoubleBottomRightCorner
			} else {
				topLeft = BorderSingleTopLeftCorner
				horiz = BorderSingleHoriz
				topRight = BorderSingleTopRightCorner
				vert = BorderSingleVert
				bottomLeft = BorderSingleBottomLeftCorner
				bottomRight = BorderSingleBottomRightCorner
			}

			terminal.SetCursorPos(b.dimensions.Left, b.dimensions.Top)
			fmt.Print(topLeft)
			for i := 1; i < b.dimensions.Width-1; i++ {
				fmt.Print(horiz)
			}
			fmt.Print(topRight)
			for i := 1; i < b.dimensions.Height-1; i++ {
				terminal.SetCursorPos(b.dimensions.Left, b.dimensions.Top+i)
				fmt.Print(vert)

				if mode == RenderModeAll {
					for i := 1; i < b.dimensions.Width-1; i++ {
						fmt.Print(" ")
					}
				} else {
					terminal.SetCursorPos(b.dimensions.Left+b.dimensions.Width-1, b.dimensions.Top+i)
				}

				fmt.Print(vert)
			}
			terminal.SetCursorPos(b.dimensions.Left, b.dimensions.Top+b.dimensions.Height)
			fmt.Print(bottomLeft)
			for i := 1; i < b.dimensions.Width-1; i++ {
				fmt.Print(horiz)
			}
			fmt.Print(bottomRight)
		}
	}
}

func (b *Box) GetChildren() []Widget {
	return nil
}

func (b *Box) CaptureInput(r string) string {
	return r
}

func (b *Box) FindInFocus(me Widget, prevInFocus Widget, inFocus Widget, focus *Focus) *Focus {
	if focus == nil {
		focus = &Focus{}
	}

	if me.CanHaveFocus() {
		if focus.Me != nil && focus.Next == nil {
			// if prev was set but next wasn't set yet, this will be the first true candidate
			// for next
			focus.Next = me
		}

		if focus.First == nil {
			focus.First = me
		}

		// last is always being updated
		focus.Last = me

		// we found the item in focus, but we have to traverse the entire application in order
		// to support focus wrapping, etc.
		if me == inFocus {
			focus.Me = me
			focus.Prev = prevInFocus
		}
	}

	// we use "me", to ensure that the child "class" receives the call, not the box parent, unless
	// nothing has been overridden
	children := me.GetChildren()
	for _, child := range children {
		child.FindInFocus(child, prevInFocus, inFocus, focus)
		prevInFocus = child
	}

	return focus
}

func (b *Box) NextInFocus(inFocus Widget) Widget {
	return nil
}

func (b *Box) CanHaveFocus() bool {
	return true
}
