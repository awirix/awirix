package tui

import (
	"fmt"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/viper"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/key"
	"github.com/awirix/awirix/scraper"
	"github.com/awirix/awirix/style"
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
		title += " " + style.Faint("by "+internal.Author())
	case *scraper.Media:
		if internal.HasInfo() {
			title += " " + style.Fg(color.Blue)(icon.Info)
		}
	}

	if l.Selected() {
		title += " " + style.Fg(color.Green)(icon.CDot)
	}

	return l.clickable(title)
}

func (l *lItem) description() string {
	const noDescription = "No description"
	switch item := l.internal.(type) {
	case *extension.Extension:
		var b strings.Builder
		if item.Passport().NSFW {
			b.WriteString(style.Fg(color.Red)("NSFW"))
			b.WriteRune(' ')
		}

		about := item.Passport().About
		if about == "" {
			about = noDescription
		}

		b.WriteString(about)
		return b.String()
	case *scraper.Media:
		description := item.Description()
		if description == "" {
			return noDescription
		}

		return description
	case *scraper.Action:
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
