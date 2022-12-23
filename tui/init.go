package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/extensions/manager"
)

func (m *model) Init() tea.Cmd {
	extensions, err := manager.InstalledExtensions()
	if err != nil {
		m.error <- err
		return nil
	}

	m.extensions = extensions

	var items []list.Item
	for _, ext := range extensions {
		items = append(items, newItem(ext))
	}

	return tea.Batch(
		m.component.extensionSelect.SetItems(items),
		m.component.textInput.Focus(),
		textinput.Blink,
	)
}
