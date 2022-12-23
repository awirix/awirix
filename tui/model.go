package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/context"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/tui/bind"
)

type model struct {
	options *Options
	error   chan error
	history *stack.Stack[state]
	context *context.Context

	extensions []*extension.Extension

	current struct {
		width, height int
		state         state
		extension     *extension.Extension
		layer         *scraper.Layer
		error         error
	}

	component struct {
		extensionSelect list.Model
		textInput       textinput.Model
		searchResults   list.Model
		layers          map[string]*list.Model
	}

	status string
	keyMap *bind.KeyMap

	style struct {
		global lipgloss.Style
	}
}

func (m *model) resize(width, height int) {
	m.current.width, m.current.height = width, height

	frameX, frameY := m.style.global.GetFrameSize()

	lists := []*list.Model{
		&m.component.extensionSelect,
		&m.component.searchResults,
	}

	for _, lst := range m.component.layers {
		lists = append(lists, lst)
	}

	for _, lst := range lists {
		lst.SetSize(width-frameX, height-frameY)
	}
}
