package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/extensions/extension"
)

func (m *model) handleLoadExtension(ext *extension.Extension) tea.Cmd {
	return func() tea.Msg {
		m.status = "Loading extension"
		ext.SetContext(m.context)
		err := ext.LoadScraper(false)
		if err != nil {
			m.error <- err
			return nil
		}

		return msgExtensionLoaded(ext)
	}
}
