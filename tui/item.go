package tui

import (
	"fmt"
	"github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
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

	return "Unnamed"
}

func (l *lItem) Title() string {
	return l.FilterValue()
}

func (l *lItem) description() string {
	switch item := l.internal.(type) {
	case *extension.Extension:
		var b strings.Builder
		if item.Passport().NSFW {
			b.WriteString(style.Fg(color.Red)("NSFW"))
			b.WriteRune(' ')
		}

		about := item.Passport().About
		if about == "" {
			about = "No description"
		}

		b.WriteString(about)
		return b.String()
	case *scraper.Media:
		return item.About()
	}

	return "No description"
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
