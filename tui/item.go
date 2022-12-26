package tui

import (
	"fmt"
	"github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
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
	return zone.Mark(l.id, s)
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

	if l.Selected() {
		title += " " + style.Fg(color.Green)(icon.CDot)
	}

	return zone.Mark(l.id, title)
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
		description := item.Description
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
