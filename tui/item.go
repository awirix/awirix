package tui

import (
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/core"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/key"
	"github.com/awirix/awirix/style"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/viper"
	"strings"
)

type lItem struct {
	id       string
	selected bool
	internal any
}

func newItem(internal any) *lItem {
	generateID := func() string {
		return fmt.Sprintf("%p", internal)
	}

	return &lItem{
		id:       generateID(),
		internal: internal,
	}
}

func (l *lItem) clickable(s string) string {
	if viper.GetBool(key.TUIClickable) {
		return zone.Mark(l.id, s)
	}

	return s
}

func (l *lItem) title() string {
	stringer, ok := l.internal.(fmt.Stringer)
	if ok {
		return stringer.String()
	}

	return fmt.Sprintf("%v", l.internal)
}

func (l *lItem) Title() string {
	title := l.title()

	switch internal := l.internal.(type) {
	case *extension.Extension:
		if manager.IsFavorite(internal) {
			title += " " + style.Fg(color.Yellow)(icon.Star)
		}

		if viper.GetBool(key.TUIShowExtensionAuthor) {
			title += " " + style.Faint(fmt.Sprintf("by %s", internal.Passport().Authors))
		}
	case *core.Media:
		if internal.HasInfo() {
			title += " " + style.Fg(color.Blue)(icon.Info)
		}
	}

	if l.Selected() {
		title += " " + style.Fg(color.Green)(icon.Square)
	}

	return l.clickable(title)
}

func (l *lItem) description() string {
	const noDescription = "No description"
	switch item := l.internal.(type) {
	case *extension.Extension:
		var b strings.Builder

		about := item.Passport().About
		if about == "" {
			about = noDescription
		}

		b.WriteString(about)

		if item.Passport().NSFW {
			b.WriteRune(' ')
			b.WriteString(style.Fg(color.Red)("NSFW"))
		}

		return b.String()
	case *core.Media:
		description := item.Description()
		if description == "" {
			return noDescription
		}

		return description
	case *core.Action:
		description := item.Description
		if description == "" {
			return noDescription
		}

		return description
	default:
		return noDescription
	}
}

func (l *lItem) Description() string {
	return l.clickable(l.description())
}

func (l *lItem) FilterValue() string {
	return l.clickable(l.title())
}

func (l *lItem) Internal() any {
	return l.internal
}

func (l *lItem) Selected() bool {
	return l.selected
}

func (l *lItem) SetSelected(selected bool) {
	l.selected = selected
}
