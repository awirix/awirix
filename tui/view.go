package tui

import zone "github.com/lrstanley/bubblezone"

func (m *model) View() string {
	switch m.current.state {
	case stateExtensionSelect:
		return zone.Scan(m.style.global.Render(m.component.extensionSelect.View()))
	default:
		return ""
	}
}
