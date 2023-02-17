package tui

import (
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	configKey "github.com/awirix/awirix/key"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"github.com/samber/mo"
	"github.com/spf13/viper"
	"golang.org/x/exp/slices"
	"strings"
	"time"
)

func newTextInput(placeholder string) textinput.Model {
	t := textinput.New()
	t.CharLimit = 80
	t.Placeholder = placeholder
	t.SetCursorMode(textinput.CursorStatic)
	t.Prompt = style.New().Foreground(color.Purple).Bold(true).Render(viper.GetString(configKey.TUIPromptSymbol))
	return t
}

func (m *model) newList(title, singular, plural string, statusMessageLifetime *time.Duration) list.Model {
	delegate := list.NewDefaultDelegate()
	border := lipgloss.Border{
		Left: "█",
	}

	delegate.ShowDescription = viper.GetBool(configKey.TUIShowDescription)

	delegate.Styles.SelectedTitle = delegate.
		Styles.
		SelectedTitle.
		Border(border, false, false, false, true)

	delegate.Styles.SelectedDesc = delegate.
		Styles.
		SelectedDesc.
		Border(border, false, false, false, true)

	l := list.New(nil, delegate, 0, 0)
	l.Title = title

	if statusMessageLifetime != nil {
		l.StatusMessageLifetime = *statusMessageLifetime
	}

	l.SetStatusBarItemName(singular, plural)
	l.Styles.NoItems = m.styles.nonListGlobal.Copy().Foreground(color.Yellow)
	l.AdditionalShortHelpKeys = func() []key.Binding {
		return m.keyMap.ShortHelp()
	}

	l.AdditionalFullHelpKeys = func() (keys []key.Binding) {
		for _, k := range m.keyMap.FullHelp() {
			keys = append(keys, k...)
		}

		return
	}

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

func listSetExtensions(exts []*extension.Extension, lst *list.Model) tea.Cmd {
	var (
		regular   = make([]*extension.Extension, 0)
		favorites = make([]*extension.Extension, 0)
	)

	for _, e := range exts {
		if manager.IsFavorite(e) {
			favorites = append(favorites, e)
		} else {
			regular = append(regular, e)
		}
	}

	nameSorter := func(a, b *extension.Extension) bool {
		return a.Passport().Name <= b.Passport().Name
	}

	slices.SortFunc(regular, nameSorter)
	slices.SortFunc(favorites, nameSorter)

	return listSetItems[*extension.Extension](
		append(favorites, regular...),
		lst,
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

func listGetSelectedItem[T any](lst *list.Model) mo.Option[T] {
	item, ok := lst.SelectedItem().(*lItem)
	if !ok {
		return mo.None[T]()
	}

	internal, ok := item.Internal().(T)
	if !ok {
		return mo.None[T]()
	}

	return mo.Some(internal)
}
