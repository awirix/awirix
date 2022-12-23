package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/samber/lo"
)

type state int

//go:generate enumer -type=state -trimprefix=state -transform=kebab-case
const (
	stateExtensionSelect state = iota + 1
	stateExtensionConfig

	stateSearch
	stateSearchResults

	stateLayer
	stateLayerResults

	statePrepare
	stateStreamOrDownloadSelection
	stateStream
	stateDownload
	stateFinal

	stateLoading
	stateError
)

func (m *model) pushState(s state) tea.Cmd {
	return func() tea.Msg {
		if m.current.state == s {
			return nil
		}

		blacklist := []state{
			stateLoading,
			statePrepare,
			stateError,
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
		var cmds = make([]tea.Cmd, 0)

		m.current.state = m.history.Pop().OrElse(m.current.state)
		if m.component.textInput.Reset() {
			cmds = append(cmds, textinput.Blink)
		}

		return tea.Batch(cmds...)
	}
}
