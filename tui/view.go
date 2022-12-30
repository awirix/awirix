package tui

import (
	"github.com/charmbracelet/lipgloss"
	zone "github.com/lrstanley/bubblezone"
	"github.com/muesli/reflow/wrap"
	"strings"
)

func (m *model) View() string {
	return zone.Scan(m.styles.global.Render(m.getCurrentStateHandler().View()))
}

func (m *model) renderLines(title string, lines ...string) string {
	// TODO: make it more precise
	var b strings.Builder
	for _, line := range lines {
		b.WriteString(line)
		b.WriteRune('\n')
	}

	body := wrap.String(b.String(), m.current.dimensions.availableWidth)
	page := m.styles.titleBar.Render(title) + "\n" + body
	help := m.styles.helpStyle.Render(m.component.help.View(m.keyMap))
	height := lipgloss.Height(page + help)

	repeat := m.current.dimensions.availableHeight - height
	if repeat < 0 {
		repeat = 0
	}

	page += strings.Repeat("\n", repeat) + help

	return page
}
