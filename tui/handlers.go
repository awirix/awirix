package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
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

		if ext.Scraper().HasLayers() {
			layers := ext.Scraper().Layers()
			m.component.layers = make(map[string]*list.Model, len(layers))
			m.current.layer = layers[0]
			for _, layer := range layers {
				lst := newList(layer.Name, "media", "media")
				m.component.layers[layer.Name] = &lst
			}

			// to update layers lists
			m.resize(m.current.width, m.current.height)
		}

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

func (m *model) handleLayer(media *scraper.Media) tea.Cmd {
	return func() tea.Msg {
		layerMedia, err := m.current.layer.Function(media)

		if err != nil {
			m.error <- err
			return nil
		}

		return msgLayerDone(layerMedia)
	}
}
