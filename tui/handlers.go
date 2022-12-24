package tui

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
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
			for _, layer := range layers {
				lst := newList(layer.String(), layer.Noun().Singular(), layer.Noun().Plural())
				m.component.layers[layer.String()] = &lst
			}

			// to update layers lists
			m.resize(m.current.width, m.current.height)
		}

		return msgExtensionLoaded(ext)
	}
}

func (m *model) handleSearch(query string) tea.Cmd {
	return func() tea.Msg {
		m.status = "Searching for " + style.Fg(color.Yellow)(query)

		search := m.current.extension.Scraper().Search()
		m.component.searchResults.Title = search.String()
		m.component.searchResults.SetStatusBarItemName(search.Noun().Singular(), search.Noun().Plural())

		media, err := search.Call(query)

		if err != nil {
			m.error <- err
			return nil
		}

		return msgSearchDone(media)
	}
}

func (m *model) handleLayer(media *scraper.Media, layer *scraper.Layer) tea.Cmd {
	return func() tea.Msg {
		if media != nil {
			m.status = "Loading " + style.Fg(color.Yellow)(media.String())
		} else {
			m.status = "Loading " + style.Fg(color.Yellow)(layer.Noun().Plural())
		}

		layerMedia, err := layer.Call(media)

		if err != nil {
			m.error <- err
			return nil
		}

		return msgLayerDone(layerMedia)
	}
}
