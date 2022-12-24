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
	"golang.org/x/exp/slices"
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

func (m *model) nextLayer() *scraper.Layer {
	layers := m.current.extension.Scraper().Layers()

	if m.current.layer == nil {
		return layers[0]
	}

	index := slices.IndexFunc(layers, func(l *scraper.Layer) bool {
		return l.Name == m.current.layer.Name
	})

	if index == -1 {
		panic("current layer is not listed in the scraper")
	}

	if index == len(layers)-1 {
		return nil
	}

	return layers[index+1]
}

func (m *model) previousLayer() *scraper.Layer {
	layers := m.current.extension.Scraper().Layers()

	if m.current.layer == nil {
		return layers[0]
	}

	index := slices.IndexFunc(layers, func(l *scraper.Layer) bool {
		return l.Name == m.current.layer.Name
	})

	if index == -1 {
		panic("current layer is not listed in the scraper")
	}

	if index == 0 {
		return nil
	}

	return layers[index-1]
}