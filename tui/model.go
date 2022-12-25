package tui

import (
	"context"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/tui/bind"
	"golang.org/x/exp/slices"
)

type model struct {
	options *Options
	error   map[*context.Context]chan error
	history *stack.Stack[state]

	extensions []*extension.Extension

	current struct {
		width, height     int
		state             state
		extension         *extension.Extension
		layer             *scraper.Layer
		error             map[*context.Context]error
		media             *scraper.Media
		context           context.Context
		contextCancelFunc context.CancelFunc
	}

	component struct {
		extensionSelect  list.Model
		textInput        textinput.Model
		searchResults    list.Model
		layers           map[string]*list.Model
		streamOrDownload list.Model
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
		&m.component.streamOrDownload,
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
		return l.Title() == m.current.layer.Title()
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
		return l.Title() == m.current.layer.Title()
	})

	if index == -1 {
		panic("current layer is not listed in the scraper")
	}

	if index == 0 {
		return nil
	}

	return layers[index-1]
}

func (m *model) cancel() {
	m.current.contextCancelFunc()
	m.current.context, m.current.contextCancelFunc = context.WithCancel(context.Background())

	if m.current.extension != nil {
		m.current.extension.SetContext(m.current.context)
	}

	m.error[&m.current.context] = make(chan error)
}
