package tui

import (
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/scraper"
	"github.com/awirix/awirix/style"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"strings"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case err := <-m.errorChan:
		m.current.cancelContext()
		return m, func() tea.Msg {
			return msgError(err)
		}
	default:
	}

	currentStateHandler := m.getCurrentStateHandler()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.Quit):
			if currentStateHandler.SoftQuit {
				return m, tea.Quit
			}
		case key.Matches(msg, m.keyMap.GoBack):
			return m, currentStateHandler.Back()
		}
	}

	return currentStateHandler.Update(msg)
}

func (m *model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgError:
		m.current.error = msg
		m.current.cancelContext()
		return m, m.pushState(stateError)
	case msgExtensionLoaded:
		m.current.extension = msg
		if m.current.extension.Scraper().HasSearch() {
			return m, m.pushState(stateSearch)
		} else {
			return m, m.handleLayer(nil, m.nextLayer())
		}
	case msgSearchDone:
		var items = make([]list.Item, len(msg))

		for i, m := range msg {
			items[i] = newItem(m)
		}

		return m, tea.Batch(
			m.component.searchResults.SetItems(items),
			m.pushState(stateSearchResults),
		)
	case msgLayerDone:
		// tea.Sequence is broken, so use deprecated tea.Sequentially
		return m, tea.Sequentially(
			listSetItems[*scraper.Media](
				msg,
				m.component.layers[m.nextLayer().String()],
			),
			func() tea.Msg {
				return msgLayerItemsSet{}
			},
		)
	case msgLayerItemsSet:
		return m, m.pushState(stateLayer)
	case msgActionDone:
		// TODO: push final state
		return m, m.popState()
	case msgMediaInfoDone:
		// to set it to the media info viewport
		m.resize(m.current.dimensions.terminalWidth, m.current.dimensions.terminalHeight)
		return m, m.pushState(stateMediaInfo)
	case msgExtensionRemoved:
		return m, tea.Sequence(
			m.pushState(stateExtensionSelect),
			m.component.extensionSelect.NewStatusMessage(fmt.Sprintf(
				"Removed %s",
				style.Fg(color.Yellow)((*extension.Extension)(msg).String()),
			)),
		)
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.component.spinner, cmd = m.component.spinner.Update(msg)
		return m, cmd
	default:
		return m, func() tea.Msg {
			return nil
		}
	}
}

func (m *model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: copy to clipboard
	return m, nil
}

func (m *model) updateExtensionRemove(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Confirm):
			return m, tea.Batch(
				m.handleExtensionRemove(m.current.extensionToRemove),
				m.pushState(stateLoading),
			)
		}
	}

	return m, nil
}

func (m *model) updateExtensionSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := &m.component.extensionSelect

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		if thisList.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, m.keyMap.Favorite):
			ext, ok := listGetSelectedItem[*extension.Extension](thisList).Get()
			if !ok {
				goto end
			}

			err := manager.ToggleFavorite(ext)
			if err != nil {
				m.errorChan <- err
			}

			return m, tea.Batch(
				m.handleExtensionsReset(),
				m.component.extensionSelect.NewStatusMessage(fmt.Sprintf(
					"Marked %s as %s",
					style.Fg(color.Yellow)(ext.String()),
					style.Fg(color.Yellow)(func() string {
						if manager.IsFavorite(ext) {
							return "favorite"
						}

						return "not favorite"
					}()),
				)),
			)
		case key.Matches(msg, m.keyMap.ExtensionRemove):
			ext, ok := listGetSelectedItem[*extension.Extension](thisList).Get()
			if !ok {
				goto end
			}

			m.current.extensionToRemove = ext
			return m, m.pushState(stateExtensionRemove)
		case key.Matches(msg, m.keyMap.ExtensionAdd):
			return m, m.pushState(stateExtensionAdd)
		case key.Matches(msg, m.keyMap.Info):
			ext, ok := listGetSelectedItem[*extension.Extension](thisList).Get()
			if !ok {
				goto end
			}

			return m, tea.Batch(
				m.pushState(stateLoading),
				m.handleExtensionInfo(ext),
			)
		case key.Matches(msg, m.keyMap.Confirm):
			ext, ok := listGetSelectedItem[*extension.Extension](thisList).Get()
			if !ok {
				goto end
			}

			return m, tea.Batch(
				m.handleLoadExtension(ext),
				m.pushState(stateLoading),
			)
		}
	}

end:
	model, cmd := m.component.extensionSelect.Update(msg)
	m.component.extensionSelect = model
	return m, cmd
}

func (m *model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Confirm):
			query := m.component.searchInput.Value()
			if strings.TrimSpace(query) == "" {
				return m, nil
			}

			return m, tea.Batch(
				m.handleSearch(query),
				m.pushState(stateLoading),
			)
		}
	}

	model, cmd := m.component.searchInput.Update(msg)
	m.component.searchInput = model
	return m, cmd
}

func (m *model) updateSearchResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := &m.component.searchResults

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		if thisList.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
		case key.Matches(msg, m.keyMap.Reset):
			m.resetSelected()
		case key.Matches(msg, m.keyMap.Info):
			media, ok := listGetSelectedItem[*scraper.Media](thisList).Get()
			if !ok {
				goto end
			}

			if !media.HasInfo() {
				goto end
			}

			return m, tea.Batch(
				m.pushState(stateLoading),
				m.handleMediaInfo(media),
			)
		case key.Matches(msg, m.keyMap.SelectAll):
			if m.current.extension.Scraper().HasLayers() {
				goto end
			}

			m.resetSelected()
			for _, item := range thisList.Items() {
				m.toggleSelect(item.(*lItem))
			}
		case key.Matches(msg, m.keyMap.Select):
			if m.current.extension.Scraper().HasLayers() {
				goto end
			}

			item, ok := thisList.SelectedItem().(*lItem)
			if !ok {
				goto end
			}

			m.toggleSelect(item)
		case key.Matches(msg, m.keyMap.Confirm):
			media, ok := listGetSelectedItem[*scraper.Media](thisList).Get()
			if !ok {
				goto end
			}

			if m.current.extension.Scraper().HasLayers() {
				return m, tea.Batch(
					m.handleLayer(media, m.nextLayer()),
					m.pushState(stateLoading),
				)
			}

			if m.current.extension.Scraper().HasActions() {
				if len(m.selectedMedia) == 0 {
					item, ok := thisList.SelectedItem().(*lItem)
					if !ok {
						goto end
					}

					m.toggleSelect(item)
				}

				return m, m.pushState(stateActionSelect)
			}

			return m, nil
		}
	}

end:
	model, cmd := m.component.searchResults.Update(msg)
	m.component.searchResults = model
	return m, cmd
}

func (m *model) updateLayer(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := m.component.layers[m.current.layer.String()]

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		if thisList.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
		case key.Matches(msg, m.keyMap.Reset):
			m.resetSelected()
		case key.Matches(msg, m.keyMap.SelectAll):
			if m.nextLayer() != nil {
				goto end
			}

			m.resetSelected()
			for _, item := range thisList.Items() {
				m.toggleSelect(item.(*lItem))
			}
		case key.Matches(msg, m.keyMap.Info):
			media, ok := listGetSelectedItem[*scraper.Media](thisList).Get()
			if !ok {
				goto end
			}

			if !media.HasInfo() {
				goto end
			}

			return m, tea.Batch(
				m.pushState(stateLoading),
				m.handleMediaInfo(media),
			)
		case key.Matches(msg, m.keyMap.Select):
			if m.nextLayer() != nil {
				goto end
			}

			item, ok := thisList.SelectedItem().(*lItem)
			if !ok {
				goto end
			}

			m.toggleSelect(item)
		case key.Matches(msg, m.keyMap.Confirm):
			if m.nextLayer() == nil {
				if !m.current.extension.Scraper().HasActions() {
					goto end
				}

				if len(m.selectedMedia) == 0 {
					item, ok := thisList.SelectedItem().(*lItem)
					if !ok {
						goto end
					}

					m.toggleSelect(item)
				}

				return m, m.pushState(stateActionSelect)
			}

			media, ok := listGetSelectedItem[*scraper.Media](thisList).Get()
			if !ok {
				goto end
			}

			return m, tea.Batch(
				m.handleLayer(media, m.nextLayer()),
				m.pushState(stateLoading),
			)
		}
	}

end:
	model, cmd := thisList.Update(msg)
	m.component.layers[m.current.layer.String()] = &model
	return m, cmd
}

func (m *model) updateActionSelect(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := &m.component.actionSelect

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		if thisList.FilterState() == list.Filtering {
			goto end
		}

		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
		case key.Matches(msg, m.keyMap.Confirm):
			action, ok := listGetSelectedItem[*scraper.Action](thisList).Get()
			if !ok {
				goto end
			}

			return m, tea.Batch(
				m.handleAction(action),
				m.pushState(stateLoading),
			)
		}
	}

end:
	model, cmd := thisList.Update(msg)
	m.component.actionSelect = model
	return m, cmd
}

func (m *model) updateMediaInfo(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.component.mediaInfo, cmd = m.component.mediaInfo.Update(msg)
	return m, cmd
}

func (m *model) updateExtensionAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO
	var cmd tea.Cmd
	m.component.extensionAddInput, cmd = m.component.extensionAddInput.Update(msg)
	return m, cmd
}
