package tui

import (
	"fmt"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/style"
	"github.com/awirix/awirix/text"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *model) Init() tea.Cmd {
	extensions, err := manager.Installed()
	if err != nil {
		m.errorChan <- err
		return nil
	}

	m.extensions = extensions

	versionMsg := fmt.Sprintf(
		"%s %s %s",
		style.Fg(color.Pink)(icon.Heart),
		text.Capitalize(app.Name),
		style.Bold(app.Version.String()),
	)

	return tea.Batch(
		listSetExtensions(extensions, &m.component.extensionSelect),
		m.component.searchInput.Focus(),
		m.component.extensionAddInput.Focus(),
		m.component.extensionSelect.NewStatusMessage(versionMsg),
	)
}
