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
		state     state
		extension *extension.Extension
		media     *scraper.Media
		layer     *scraper.Layer
		query     string
		error     error
	}

	component struct {
		extensionSelect list.Model
		textInput       textinput.Model
		searchResults   list.Model
		layerResults    list.Model
	}

	status string
	keyMap *bind.KeyMap

	style struct {
		global lipgloss.Style
	}
}

func (m *model) resize(width, height int) {
	frameX, frameY := m.style.global.GetFrameSize()

	for _, lst := range []*list.Model{
		&m.component.extensionSelect,
		&m.component.searchResults,
	} {
		lst.SetSize(width-frameX, height-frameY)
	}
}
