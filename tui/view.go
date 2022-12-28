package tui

import "strings"

func (m *model) View() string {
	return m.getCurrentStateHandler().View()
}

func (m *model) renderLines(title string, lines ...string) string {
	l := m.style.titleBar.Render(title) + "\n\n" + strings.Join(lines, "\n")
	return m.style.global.Render(l)
}
