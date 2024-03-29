package tui

import (
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/style"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
	"strings"
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

	stateExtensionRemove
	stateExtensionAdd
)

type handler struct {
	Update   func(msg tea.Msg) (tea.Model, tea.Cmd)
	View     func() string
	Back     func() tea.Cmd
	SoftQuit bool
}

func (m *model) getCurrentStateHandler() *handler {
	// TODO: optimize (create this map only once and reuse it)

	defaultBack := func() tea.Cmd {
		if m.current.cancelContext != nil {
			m.current.cancelContext()
		}
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
					"Please, wait",
					m.component.spinner.View()+style.Faint(m.text.status),
				)
			},
			Back: defaultBack,
		},

		stateError: {
			Update: m.updateError,
			View: func() string {
				var err string
				if m.current.error == nil {
					err = "Unknown error occurred... something went wrong!"
				} else {
					err = m.current.error.Error()
				}
				err = strings.TrimSpace(err)
				return m.renderLines(
					m.styles.titleError.Render("️Error"),
					"An error occurred",
					style.Fg(color.Red)(err),
				)
			},
			Back:     defaultBack,
			SoftQuit: true,
		},

		stateExtensionSelect: {
			Update: m.updateExtensionSelect,
			View: func() string {
				return m.component.extensionSelect.View()
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.extensionSelect)
			},
		},

		stateSearch: {
			Update: m.updateSearch,
			View: func() string {
				return m.renderLines(
					m.styles.title.Render(m.text.searchTitle),
					"Enter a search query",
					m.component.searchInput.View(),
				)
			},
			Back: defaultBack,
		},

		stateSearchResults: {
			Update: m.updateSearchResults,
			View: func() string {
				return m.component.searchResults.View()
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.searchResults)
			},
		},

		stateLayer: {
			Update: m.updateLayer,
			View: func() string {
				current := m.component.layers[m.current.layer.String()]
				return current.View()
			},
			Back: func() tea.Cmd {
				return listBack(m.component.layers[m.current.layer.String()])
			},
		},

		stateActionSelect: {
			Update: m.updateActionSelect,
			View: func() string {
				return m.component.actionSelect.View()
			},
			Back: func() tea.Cmd {
				return listBack(&m.component.actionSelect)
			},
		},

		stateMediaInfo: {
			Update: m.updateMediaInfo,
			View: func() string {
				return m.renderLines(
					m.styles.title.Render(m.text.mediaInfoTitle),
					m.text.mediaInfoName,
					m.component.mediaInfo.View(),
				)
			},
			Back:     defaultBack,
			SoftQuit: true,
		},

		stateExtensionRemove: {
			Update: m.updateExtensionRemove,
			View: func() string {
				return m.renderLines(
					m.styles.title.Render("Remove"),
					"Confirm removal of extension",
					fmt.Sprintf(`Are you sure you want to %s the extension %s`, style.Fg(color.Red)("remove"), style.Fg(color.Yellow)(m.current.extensionToRemove.String())),
					"",
					fmt.Sprintf(
						"%s   %s",
						"No "+style.Faint(m.keyMap.GoBack.Help().Key),
						"Yes "+style.Faint(m.keyMap.Confirm.Help().Key),
					),
				)
			},
			Back:     defaultBack,
			SoftQuit: true,
		},

		stateExtensionAdd: {
			Update: m.updateExtensionAdd,
			View: func() string {
				// TODO
				return m.renderLines(
					m.styles.title.Render("Add"),
					"Install a new extension",
					"NOT IMPLEMENTED",
					m.component.extensionAddInput.View(),
				)
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
			stateExtensionRemove,
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
			m.component.searchInput.Reset()
		case stateSearchResults:
			m.component.searchResults.ResetSelected()
			m.resetSelected()
		case stateLayer:
			m.resetSelected()
		}

		m.component.mediaInfo.GotoTop()
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
