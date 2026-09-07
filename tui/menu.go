package tui

import (
	"fmt"
	"sync"

	"github.com/exadrift/go/ansi/style"
)

type Menu struct {
	*Box
	contents      []string
	index         map[string]int
	selectedIndex int
	selectHandler func(int, string) any
	completer     func(any)
	busyLabel     string

	defaultStyles  style.Styles
	selectedStyles style.Styles
}

func NewMenu(contents ...string) *Menu {
	menu := &Menu{
		Box: NewBox(),
	}

	menu.SetSelectedStyles(style.Blue.Bg(), style.White.Fg())
	menu.SetContents(contents...)

	return menu
}

// SetSelectHandler allows you to define a callback which will be executed when selection changes
func (m *Menu) SetSelectHandler(h func(selectedIndex int, selectedItem string) any, options ...*Option) *Menu {
	m.selectHandler = h

	for _, option := range options {
		switch option.optionType {
		case BusyModal:
			busyModalData := option.data.(*BusyModalData)
			m.busyLabel = busyModalData.Label
			m.completer = busyModalData.Completer
		default:
			panic(fmt.Errorf("unknown menu option: %v", option.optionType))
		}
	}

	return m
}

func (m *Menu) SetDefaultStyles(styles ...style.Style) *Menu {
	m.defaultStyles = styles[:]
	return m
}

func (m *Menu) SetSelectedStyles(styles ...style.Style) *Menu {
	m.selectedStyles = styles[:]
	return m
}

func (m *Menu) SetContents(contents ...string) *Menu {
	m.contents = make([]string, len(contents))
	m.index = make(map[string]int, len(contents))
	m.selectedIndex = 0

	m.ResetScrollPosition()

	for i, item := range contents {
		// sorry, menus shouldn't have any ANSI codes in them
		m.contents[i] = style.StripAnsi(item)
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

	return m
}

func (m *Menu) SetSelectedItem(item string) *Menu {
	index := m.index[item]
	m.SetSelectedIndex(index)

	return m
}

func (m *Menu) Render(mode RenderMode, focusItem Widget) {
	contentLength := len(m.contents)
	m.EnableScrollHandle(contentLength > m.contentDimensions.Height)

	m.RenderWithScroll(mode, focusItem, contentLength, m.selectedIndex, func(index int) string {
		menuLabel := Pad(Constrain(m.contents[index], m.contentDimensions.Width), m.contentDimensions.Width)
		switch index {
		case m.selectedIndex:
			return fmt.Sprintf("%s%s%s", m.selectedStyles.Ansi(), menuLabel, StyleReset)
		default:
			return fmt.Sprintf("%s%s%s", m.defaultStyles.Ansi(), menuLabel, StyleReset)
		}
	})
}

func (m *Menu) CaptureInput(r string) string {
	switch r {
	case appSingleton.keyBindings.SelectionPrev:
		m.SetSelectedIndex(m.selectedIndex - 1)
	case appSingleton.keyBindings.SelectionNext:
		m.SetSelectedIndex(m.selectedIndex + 1)
	case appSingleton.keyBindings.Trigger:
		if m.selectHandler != nil {
			if m.completer == nil {
				m.selectHandler(m.selectedIndex, m.contents[m.selectedIndex])
				return RenderFullCode
			}

			appSingleton.ShowLoader(m.busyLabel)
			go func() {
				// this is an async operation
				wg := sync.WaitGroup{}
				wg.Add(1)
				var resp any
				go func() {
					defer wg.Done()
					resp = m.selectHandler(m.selectedIndex, m.contents[m.selectedIndex])
				}()
				wg.Wait()
				appSingleton.Async(func() {
					m.completer(resp)
					appSingleton.renderAll()
					appSingleton.HideLoader()
				})
			}()
		}
	default:
		return r
	}

	return ""
}
