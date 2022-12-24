package tui

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/style"
)

type state int

//go:generate enumer -type=state -trimprefix=state -transform=kebab-case
const (
	stateExtensionSelect state = iota + 1
	stateExtensionConfig

	stateSearch
	stateSearchResults

	stateLayer

	statePrepare
	stateStreamOrDownloadSelection
	stateStream
	stateDownload
	stateFinal

	stateLoading
	stateError
)

type handler struct {
	Update func(msg tea.Msg) (tea.Model, tea.Cmd)
	View   func() string
}

func (m *model) getCurrentStateHandler() *handler {
	// TODO: optimize (create this map only once and reuse it)

	h, ok := map[state]*handler{
		stateLoading: {
			Update: m.updateLoading,
			View:   func() string { return m.status },
		},

		stateError: {
			Update: m.updateError,
			View: func() string {
				return style.Fg(color.Red)(m.current.error.Error())
			},
		},

		stateExtensionSelect: {
			Update: m.updateExtensionSelect,
			View: func() string {
				return zone.Scan(m.style.global.Render(m.component.extensionSelect.View()))
			},
		},

		stateSearch: {
			Update: m.updateSearch,
			View:   m.component.textInput.View,
		},

		stateSearchResults: {
			Update: m.updateSearchResults,
			View: func() string {
				return zone.Scan(m.style.global.Render(m.component.searchResults.View()))
			},
		},

		stateLayer: {
			Update: m.updateLayer,
			View: func() string {
				current := m.component.layers[m.current.layer.Name]
				return zone.Scan(m.style.global.Render(current.View()))
			},
		},
	}[m.current.state]

	if !ok {
		panic(fmt.Sprintf(`Unknown state "%s"`, m.current.state.String()))
	}

	return h
}

func (m *model) pushState(s state) tea.Cmd {
	return func() tea.Msg {
		blacklist := []state{
			stateLoading,
			statePrepare,
			stateError,
		}

		// Layers are special case
		// Basically, when we receive `pushState(stateLayer)`
		// we just set current layer to the next in a sequence.
		if s == stateLayer {
			if nextLayer := m.nextLayer(); nextLayer != nil {
				m.current.layer = m.nextLayer()
			}
		}

		if m.current.state == s {
			return nil
		}

		if !lo.Contains[state](blacklist, m.current.state) {
			m.history.Push(m.current.state)
		}

		m.current.state = s
		return nil
	}
}

func (m *model) popState() tea.Cmd {
	return func() tea.Msg {
		if m.current.state == stateLayer {
			previous := m.previousLayer()

			m.component.layers[m.current.layer.Name].ResetSelected()

			// if we're going back from the first layer
			if previous == nil {
				// reset layers lists
				m.component.layers = make(map[string]*list.Model)
				m.current.layer = nil
				goto regular
			}

			m.current.layer = previous
		}

	regular:
		var cmds = make([]tea.Cmd, 0)

		popped, ok := m.history.Pop().Get()
		if !ok {
			return nil
		}

		switch m.current.state {
		case stateSearch:
			if m.component.textInput.Reset() {
				cmds = append(cmds, textinput.Blink)
			}
		case stateSearchResults:
			m.component.searchResults.ResetSelected()
		}

		m.current.state = popped
		return tea.Batch(cmds...)
	}
}
