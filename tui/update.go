package tui

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case err := <-m.error[&m.current.context]:
		m.current.error[&m.current.context] = err
		m.cancel()
		return m, m.pushState(stateError)
	default:
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.GoBack):
			m.cancel()
			return m, m.popState()
		}
	}

	return m.getCurrentStateHandler().Update(msg)
}

func (m *model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
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
				m.component.layers[m.nextLayer().Title()],
			),
			func() tea.Msg {
				return msgLayerItemsSet{}
			},
		)
	case msgLayerItemsSet:
		return m, m.pushState(stateLayer)
	case msgPrepareDone:
		m.current.media = msg
		return m, m.pushState(stateStreamOrDownloadSelection)
	default:
		// trigger update
		return m, func() tea.Msg { return msg }
	}
}

func (m *model) updateError(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Quit):
			return m, tea.Quit
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
		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
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
			query := m.component.textInput.Value()
			if query == "" {
				return m, nil
			}

			return m, tea.Batch(
				m.handleSearch(query),
				m.pushState(stateLoading),
			)
		}
	}

	model, cmd := m.component.textInput.Update(msg)
	m.component.textInput = model
	return m, cmd
}

func (m *model) updateSearchResults(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := &m.component.searchResults

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
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

			return m, tea.Batch(
				m.handlePrepare(media),
				m.pushState(stateLoading),
			)
		}
	}

end:
	model, cmd := m.component.searchResults.Update(msg)
	m.component.searchResults = model
	return m, cmd
}

func (m *model) updateLayer(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := m.component.layers[m.current.layer.Title()]

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
		case key.Matches(msg, m.keyMap.Confirm):
			media, ok := listGetSelectedItem[*scraper.Media](thisList).Get()
			if !ok {
				goto end
			}

			if m.nextLayer() == nil {
				return m, tea.Batch(
					m.handlePrepare(media),
					m.pushState(stateLoading),
				)
			}

			return m, tea.Batch(
				m.handleLayer(media, m.nextLayer()),
				m.pushState(stateLoading),
			)
		}
	}

end:
	model, cmd := thisList.Update(msg)
	m.component.layers[m.current.layer.Title()] = &model
	return m, cmd
}

func (m *model) updateStreamOrDownload(msg tea.Msg) (tea.Model, tea.Cmd) {
	thisList := &m.component.streamOrDownload

	switch msg := msg.(type) {
	case tea.MouseMsg:
		listHandleMouseMsg(msg, thisList)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Reverse):
			return m, listReverseItems(thisList)
		case key.Matches(msg, m.keyMap.Confirm):
			item, ok := listGetSelectedItem[variant](thisList).Get()
			if !ok {
				goto end
			}

			switch item {
			case variantStream:
				// TODO
			case variantDownload:
				// TODO
			}
		}
	}

end:
	model, cmd := thisList.Update(msg)
	m.component.streamOrDownload = model
	return m, cmd
}

func (m *model) updateStream(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO

	return m, nil
}

func (m *model) updateDownload(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO

	return m, nil
}
