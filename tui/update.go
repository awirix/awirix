package tui

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/tui/bind"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, bind.ForceQuit):
			return m, tea.Quit
		}
	}

	switch m.current.state {
	case stateExtensionSelect:
		return m.updateExtensionSelect(msg)
	default:
		return m, nil
	}
}

func (m *model) updateExtensionSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.MouseMsg:
		switch msg.Type {
		case tea.MouseWheelUp:
			m.component.extensionSelect.CursorUp()
		case tea.MouseWheelDown:
			m.component.extensionSelect.CursorDown()
		case tea.MouseLeft:
			for i, listItem := range m.component.extensionSelect.VisibleItems() {
				item, _ := listItem.(*lItem)
				if zone.Get(item.id).InBounds(msg) {
					m.component.extensionSelect.Select(i)
					break
				}
			}
		}

		return m, nil
	}

	model, cmd := m.component.extensionSelect.Update(msg)
	m.component.extensionSelect = model
	return m, cmd
}
