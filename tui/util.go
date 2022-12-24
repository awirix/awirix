package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/option"
	"strings"
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

func listReverseItems(lst *list.Model) tea.Cmd {
	var b strings.Builder
	_, _ = log.WriteSuccessf(&b, "Reversed")
	items := lst.Items()
	return tea.Batch(
		lst.SetItems(lo.Reverse(items)),
		lst.NewStatusMessage(b.String()),
	)
}

func listSetItems[T any](items []T, lst *list.Model) tea.Cmd {
	var listItems = make([]list.Item, len(items))

	for i, m := range items {
		listItems[i] = newItem(m)
	}

	return lst.SetItems(listItems)
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
