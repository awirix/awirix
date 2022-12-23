package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/extensions/extension"
)

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	select {
	case err := <-m.error:
		m.current.error = err
		m.context.Cancel()
		m.context.Reset()
		return m, m.pushState(stateError)
	default:
		goto main
	}

main:
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.resize(msg.Width, msg.Height)
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.ForceQuit):
			return m, tea.Quit
		case key.Matches(msg, m.keyMap.GoBack):
			return m, m.popState()
		}
	}

	switch m.current.state {
	case stateLoading:
		return m.updateLoading(msg)
	case stateError:
		return m.updateError(msg)
	case stateExtensionSelect:
		return m.updateExtensionSelect(msg)
	case stateSearch:
		return m.updateSearch(msg)
	case stateSearchResults:
		return m.updateSearchResults(msg)
	default:
		panic(fmt.Sprintf(`Unknown state "%s"`, m.current.state.String()))
	}
}

func (m *model) updateLoading(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgExtensionLoaded:
		m.current.extension = msg
		if m.current.extension.Scraper().HasSearch() {
			return m, m.pushState(stateSearch)
		} else {
			return m, m.pushState(stateLayer)
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
		switch msg.Type {
		case tea.MouseWheelUp:
			thisList.CursorUp()
		case tea.MouseWheelDown:
			thisList.CursorDown()
		case tea.MouseLeft:
			for i, listItem := range thisList.VisibleItems() {
				item, _ := listItem.(*lItem)
				if zone.Get(item.id).InBounds(msg) {
					thisList.Select(i)
					goto end
				}
			}
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Confirm):
			item, ok := thisList.SelectedItem().(*lItem)
			if !ok {
				goto end
			}

			ext, ok := item.Internal().(*extension.Extension)
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keyMap.Confirm):
			// TODO
			goto end
		}
	}

end:
	model, cmd := m.component.searchResults.Update(msg)
	m.component.searchResults = model
	return m, cmd
}
