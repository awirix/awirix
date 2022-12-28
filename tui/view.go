package tui

import "strings"

func (m *model) View() string {
	return m.getCurrentStateHandler().View()
}

func (m *model) renderLines(title string, lines ...string) string {
	l := m.styles.titleBar.Render(title) + "\n" + strings.Join(lines, "\n")
	return m.styles.global.Render(l)
}
