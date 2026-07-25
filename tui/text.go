package tui

import "github.com/exadrift/go/tui/internal/terminal"

type Text struct {
	*Box
	Contents string
}

func NewText(contents string) *Text {
	return &Text{
		Box:      NewBox(),
		Contents: contents,
	}
}

func (t *Text) Render(mode RenderMode, focusItem Widget) {
	t.Box.Render(mode, focusItem)

	dimensions := t.GetContentDimensions()
	terminal.SetCursorPos(dimensions.Left, dimensions.Top)
}
