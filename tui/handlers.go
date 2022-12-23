package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/extensions/extension"
)

func (m *model) handleLoadExtension(ext *extension.Extension) tea.Cmd {
	return func() tea.Msg {
		m.status = "Loading extension"
		ext.SetContext(m.context)

		if !ext.IsScraperLoaded() {
			err := ext.LoadScraper(false)
			if err != nil {
				m.error <- err
				return nil
			}
		}

		ext.Scraper().SetProgress(func(message string) {
			m.status = message
		})

		return msgExtensionLoaded(ext)
	}
}

func (m *model) handleSearch(query string) tea.Cmd {
	return func() tea.Msg {
		media, err := m.current.extension.Scraper().Search(query)

		if err != nil {
			m.error <- err
			return nil
		}

		return msgSearchDone(media)
	}
}
