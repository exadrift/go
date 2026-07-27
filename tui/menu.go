package tui

import (
	"fmt"

	"github.com/exadrift/go/tui/internal/terminal"
)

type Menu struct {
	*Box
	contents       []string
	index          map[string]int
	selectedIndex  int
	scrollPosition int
	selectHandler  func(int, string)
}

func NewMenu(contents ...string) *Menu {
	menu := &Menu{
		Box: NewBox(),
	}

	menu.SetContents(contents...)

	return menu
}

func (m *Menu) SetSelectHandler(h func(selectedIndex int, selectedItem string)) *Menu {
	m.selectHandler = h
	return m
}

func (m *Menu) SetContents(contents ...string) *Menu {
	m.contents = make([]string, len(contents))
	m.index = make(map[string]int, len(contents))
	m.selectedIndex = 0
	m.scrollPosition = 0
	for i, item := range contents {
		// sorry, menus shouldn't have any ANSI codes in them
		m.contents[i] = StripAnsiCodes(item)
	}
	for i, val := range contents {
		m.index[val] = i
	}

	return m
}

func (m *Menu) SetSelectedIndex(index int) *Menu {
	if index < 0 {
		index = len(m.contents) - 1
	} else if index > len(m.contents)-1 {
		index = 0
	}
	m.selectedIndex = index

	contentDims := m.GetContentDimensions()
	if m.selectedIndex > contentDims.Height-1+m.scrollPosition {
		m.scrollPosition = m.selectedIndex - (contentDims.Height - 1)
	} else if m.selectedIndex < m.scrollPosition {
		m.scrollPosition = m.selectedIndex
	}

	return m
}

func (m *Menu) SetSelectedItem(item string) *Menu {
	index := m.index[item]
	m.SetSelectedIndex(index)

	return m
}

func (m *Menu) Render(mode RenderMode, focusItem Widget) {
	m.Box.Render(mode, focusItem)

	dimensions := m.GetContentDimensions()
	curRow := 0
	for i := m.scrollPosition; i < len(m.contents); i++ {
		if curRow >= dimensions.Height {
			break
		}
		terminal.SetCursorPos(dimensions.Left, dimensions.Top+curRow)
		menuLabel := Pad(Constrain(m.contents[i], dimensions.Width), dimensions.Width)
		if i == m.selectedIndex {
			fmt.Print(StyleFgBg(White, Blue, menuLabel))
		} else {
			fmt.Print(StyleFg(White, menuLabel))
		}

		curRow++
	}
}

func (m *Menu) CaptureInput(r string) string {
	switch r {
	case UpArrow:
		m.SetSelectedIndex(m.selectedIndex - 1)
	case DownArrow:
		m.SetSelectedIndex(m.selectedIndex + 1)
	case Enter:
		if m.selectHandler != nil {
			m.selectHandler(m.selectedIndex, m.contents[m.selectedIndex])
			return RenderFullCode
		}
	default:
		return r
	}

	return ""
}
