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

	stateSearch
	stateSearchResults

	stateLayer

	stateActionSelect

	stateMediaInfo

	stateLoading
	stateError
)

type handler struct {
	Update func(msg tea.Msg) (tea.Model, tea.Cmd)
	View   func() string
	Back   func() tea.Cmd
}

func (m *model) getCurrentStateHandler() *handler {
	// TODO: optimize (create this map only once and reuse it)

	defaultBack := func() tea.Cmd {
		m.current.cancelContext()
		return m.popState()
	}

	listBack := func(l *list.Model) tea.Cmd {

		if l.FilterState() != list.Unfiltered {
			l.ResetFilter()
			return nil
		}

		return defaultBack()
	}

	h, ok := map[state]*handler{
		stateLoading: {
			Update: m.updateLoading,
			View: func() string {
				return m.renderLines(
					m.styles.title.Render("Loading"),
					m.component.spinner.View()+style.Faint(m.status),
				)
			},
			Back: defaultBack,
		},

		stateError: {
			Update: m.updateError,
			View: func() string {
				return m.renderLines(
					m.styles.titleError.Render("Error"),
					style.Fg(color.Red)(m.current.error.Error()),
				)
			},
			Back: defaultBack,
		},

		stateExtensionSelect: {
			Update: m.updateExtensionSelect,
			View: func() string {
				return zone.Scan(m.styles.global.Render(m.component.extensionSelect.View()))
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.extensionSelect)
			},
		},

		stateSearch: {
			Update: m.updateSearch,
			View: func() string {
				return m.renderLines(
					// TODO: use title from extension
					m.styles.title.Render("Search"),
					m.component.textInput.View(),
				)
			},
			Back: defaultBack,
		},

		stateSearchResults: {
			Update: m.updateSearchResults,
			View: func() string {
				return zone.Scan(m.styles.global.Render(m.component.searchResults.View()))
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.searchResults)
			},
		},

		stateLayer: {
			Update: m.updateLayer,
			View: func() string {
				current := m.component.layers[m.current.layer.String()]
				return zone.Scan(m.styles.global.Render(current.View()))
			},
			Back: func() tea.Cmd {
				return listBack(m.component.layers[m.current.layer.String()])
			},
		},

		stateActionSelect: {
			Update: m.updateActionSelect,
			View: func() string {
				return zone.Scan(m.styles.global.Render(m.component.actionSelect.View()))
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.actionSelect)
			},
		},

		stateMediaInfo: {
			Update: m.updateMediaInfo,
			View: func() string {
				//view := wrap.String(m.component.mediaInfo.View(), m.component.mediaInfo.Width)
				return m.styles.global.Render(m.component.mediaInfo.View())
			},
			Back: defaultBack,
		},
	}[m.current.state]

	if !ok {
		panic(fmt.Sprintf(`Unknown state "%s"`, m.current.state.String()))
	}

	return h
}

func (m *model) pushState(s state) tea.Cmd {
	return func() tea.Msg {
		// states that should not be pushed to the history
		// since it makes no sense to return to them
		blacklist := []state{
			stateLoading,
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

		if s == stateLoading {
			m.resetSpinner()
			return m.component.spinner.Tick()
		}

		if s == stateSearch {
			return textinput.Blink
		}

		return nil
	}
}

func (m *model) popState() tea.Cmd {
	return func() tea.Msg {
		if m.current.state == stateLayer {
			m.component.layers[m.current.layer.String()].ResetSelected()
			m.current.layer = m.previousLayer()
		}

		var cmds = make([]tea.Cmd, 0)

		popped, ok := m.history.Pop().Get()
		if !ok {
			return nil
		}

		switch m.current.state {
		case stateSearch:
			m.component.textInput.Reset()
		case stateSearchResults:
			m.component.searchResults.ResetSelected()
			m.resetSelected()
		case stateLayer:
			m.resetSelected()
		}

		cmds = append(cmds, m.resetListStatusMessages())
		m.current.state = popped
		return tea.Batch(cmds...)
	}
}

func (m *model) resetListStatusMessages() tea.Cmd {
	var cmds []tea.Cmd
	for _, l := range m.lists() {
		cmds = append(cmds, l.NewStatusMessage(""))
	}

	return tea.Batch(cmds...)
}

func (m *model) resetSelected() {
	for i := range m.selectedMedia {
		i.SetSelected(false)
	}
	m.selectedMedia = make(map[*lItem]struct{})
}
