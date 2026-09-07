package tui

import "github.com/exadrift/go/ansi/style"

type Text struct {
	*Box
	Contents style.Text
}

func NewText(contents any) *Text {
	return &Text{
		Box:      NewBox(),
		Contents: ProcessToStyledText(contents),
	}
}

func (t *Text) CaptureInput(r string) string {
	switch r {
	case appSingleton.keyBindings.ScrollUp:
		t.scrollWindow.ScrollUp()
	case appSingleton.keyBindings.ScrollDown:
		t.scrollWindow.ScrollDown()
	default:
		return r
	}

	return ""
}

func (t *Text) Render(mode RenderMode, focusItem Widget) {
	dimensions := t.GetContentDimensions()
	t.Contents.Render()

	lines := t.Contents.Render(style.WithWidthConstraint(dimensions.Width), style.WithDefaultStyles(t.defaultStyles...))
	t.scrollWindow.scrollHandleEnabled = len(lines) > dimensions.Height

	t.RenderWithScroll(mode, focusItem, len(lines), -1, func(index int) string {
		return Pad(lines[index], dimensions.Width)
	})
}
