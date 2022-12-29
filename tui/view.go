package tui

import (
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func (m *model) View() string {
	return m.getCurrentStateHandler().View()
}

func (m *model) renderLines(title string, lines ...string) string {
	l := m.styles.titleBar.Render(title) + "\n" + strings.Join(lines, "\n")
	help := m.styles.helpStyle.Render(m.component.help.View(m.keyMap))
	height := lipgloss.Height(l + help)
	l += strings.Repeat("\n", m.current.height-height) + help

	// TODO: add help to the bottom
	return m.styles.global.Render(l)
}
