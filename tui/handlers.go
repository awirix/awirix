package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
)

func (m *model) handleWrapper(cmd tea.Cmd) tea.Cmd {
	return func() tea.Msg {
		m.resetContext()
		return cmd()
	}
}

func (m *model) handleLoadExtension(ext *extension.Extension) tea.Cmd {
	return m.handleWrapper(func() tea.Msg {
		m.text.status = "Loading extension"

		if !ext.IsScraperLoaded() {
			err := ext.LoadScraper(false)
			if err != nil {
				return msgError(err)
			}
		}

		if ext.Scraper().HasSearch() {
			search := ext.Scraper().Search()
			m.component.searchResults.Title = search.Subtitle()
			m.component.textInput.Placeholder = search.Placeholder()
			m.component.searchResults.SetStatusBarItemName(search.Noun.Singular(), search.Noun.Plural())
			m.text.searchTitle = search.String()
		}

		if ext.Scraper().HasLayers() {
			layers := ext.Scraper().Layers()
			m.component.layers = make(map[string]*list.Model, len(layers))
			for i, layer := range layers {
				lst := m.newList(
					fmt.Sprintf("%s - %d/%d", layer.String(), i+1, len(layers)),
					layer.Noun.Singular(),
					layer.Noun.Plural(),
				)
				m.component.layers[layer.String()] = &lst
			}

			// to update layers lists
			m.resize(m.current.dimensions.terminalWidth, m.current.dimensions.terminalHeight)
		}

		// it returns tea cmd it's just nil if we don't have any filter applied,
		// so we can safely ignore it here
		listSetItems[*scraper.Action](ext.Scraper().Actions(), &m.component.actionSelect)

		return msgExtensionLoaded(ext)
	})
}

func (m *model) handleSearch(query string) tea.Cmd {
	return m.handleWrapper(func() tea.Msg {
		m.text.status = "Searching for " + style.Fg(color.Yellow)(query)

		search := m.current.extension.Scraper().Search()
		media, err := search.Call(query)

		if err != nil {
			//m.error <- err
			return msgError(err)
		}

		return msgSearchDone(media)
	})
}

func (m *model) handleLayer(media *scraper.Media, layer *scraper.Layer) tea.Cmd {
	return m.handleWrapper(func() tea.Msg {
		if media != nil {
			m.text.status = "Loading " + style.Fg(color.Yellow)(media.String())
		} else {
			m.text.status = "Loading " + style.Fg(color.Yellow)(layer.Noun.Plural())
		}

		layerMedia, err := layer.Call(media)

		if err != nil {
			//m.error <- err
			return msgError(err)
		}

		return msgLayerDone(layerMedia)
	})
}

func (m *model) handleAction(action *scraper.Action) tea.Cmd {
	return m.handleWrapper(func() tea.Msg {
		m.text.status = "Performing " + style.Fg(color.Yellow)(action.String()) + " action"

		var media = make([]*scraper.Media, 0)
		for item := range m.selectedMedia {
			media = append(media, item.Internal().(*scraper.Media))
		}

		err := action.Call(media)
		if err != nil {
			//m.error <- err
			return msgError(err)
		}

		return msgActionDone(action)
	})
}

func (m *model) handleMediaInfo(media *scraper.Media) tea.Cmd {
	return m.handleWrapper(func() tea.Msg {
		m.text.status = "Loading info for " + style.Fg(color.Yellow)(media.String())
		info, err := media.Info()
		if err != nil {
			return msgError(err)
		}

		m.text.mediaInfoName = media.String()
		m.current.mediaInfo = info
		return msgMediaInfoDone{}
	})
}
