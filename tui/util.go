package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/option"
)

func newList(title, singular, plural string) list.Model {
	delegate := list.NewDefaultDelegate()
	border := lipgloss.Border{
		Left: "█",
	}

	delegate.Styles.SelectedTitle = delegate.
		Styles.
		SelectedTitle.
		Bold(true).
		Border(border, false, false, false, true)

	delegate.Styles.SelectedDesc = delegate.
		Styles.
		SelectedDesc.
		Border(border, false, false, false, true)

	l := list.New(nil, delegate, 0, 0)
	l.Title = title
	l.SetStatusBarItemName(singular, plural)
	return l
}

func listHandleMouseMsg(msg tea.MouseMsg, lst *list.Model) {
	switch msg.Type {
	case tea.MouseWheelUp:
		lst.CursorUp()
	case tea.MouseWheelDown:
		lst.CursorDown()
	case tea.MouseLeft:
		for i, listItem := range lst.VisibleItems() {
			item, _ := listItem.(*lItem)
			if zone.Get(item.id).InBounds(msg) {
				lst.Select(i)
				break
			}
		}
	}
}

func listGetSelectedItem[T any](lst *list.Model) *option.Option[T] {
	item, ok := lst.SelectedItem().(*lItem)
	if !ok {
		return option.None[T]()
	}

	internal, ok := item.Internal().(T)
	if !ok {
		return option.None[T]()
	}

	return option.Some(internal)
}
