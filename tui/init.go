package tui

import (
	"fmt"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/style"
	"github.com/awirix/awirix/text"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	extensions, err := manager.InstalledExtensions()
	if err != nil {
		m.errorChan <- err
		return nil
	}

	m.extensions = extensions

	var items []list.Item
	for _, ext := range extensions {
		items = append(items, newItem(ext))
	}

	versionMsg := fmt.Sprintf(
		"%s %s %s",
		style.Fg(color.Pink)(icon.Heart),
		text.Capitalize(app.Name),
		style.Bold(app.Version.String()),
	)

	return tea.Batch(
		m.component.extensionSelect.SetItems(items),
		m.component.textInput.Focus(),
		m.component.extensionSelect.NewStatusMessage(versionMsg),
	)
}
