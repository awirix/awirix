package tui

import (
	"context"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/tui/bind"
	"golang.org/x/exp/slices"
)

type model struct {
	options   *Options
	errorChan chan error
	history   *stack.Stack[state]

	extensions []*extension.Extension

	// algebraic set
	selectedMedia map[*lItem]struct{}

	current struct {
		width, height int
		state         state
		extension     *extension.Extension
		layer         *scraper.Layer
		error         error
		context       context.Context
		cancelContext context.CancelFunc
	}

	component struct {
		extensionSelect list.Model
		textInput       textinput.Model
		searchResults   list.Model
		layers          map[string]*list.Model
		actionSelect    list.Model
		spinner         spinner.Model
	}

	status string
	keyMap *bind.KeyMap

	style struct {
		global,
		title, titleError, titleBar lipgloss.Style
	}
}

func (m *model) lists() []*list.Model {
	lists := []*list.Model{
		&m.component.extensionSelect,
		&m.component.searchResults,
		&m.component.actionSelect,
	}

	for _, lst := range m.component.layers {
		lists = append(lists, lst)
	}

	return lists
}

func (m *model) resize(width, height int) {
	m.current.width, m.current.height = width, height
	frameX, frameY := m.style.global.GetFrameSize()

	for _, lst := range m.lists() {
		lst.SetSize(width-frameX, height-frameY)
	}
}

func (m *model) nextLayer() *scraper.Layer {
	layers := m.current.extension.Scraper().Layers()

	if m.current.layer == nil {
		return layers[0]
	}

	index := slices.IndexFunc(layers, func(l *scraper.Layer) bool {
		return l.String() == m.current.layer.String()
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
		return l.String() == m.current.layer.String()
	})

	if index == -1 {
		panic("current layer is not listed in the scraper")
	}

	if index == 0 {
		return nil
	}

	return layers[index-1]
}

func (m *model) resetContext() {
	m.current.context, m.current.cancelContext = context.WithCancel(context.Background())
	if m.current.extension != nil {
		m.injectContext(m.current.extension)
	}
}

func (m *model) toggleSelect(item *lItem) {
	if _, ok := m.selectedMedia[item]; ok {
		delete(m.selectedMedia, item)
		item.SetSelected(false)
	} else {
		m.selectedMedia[item] = struct{}{}
		item.SetSelected(true)
	}
}

func (m *model) injectContext(ext *extension.Extension) {
	ext.SetContext(m.current.context)
	ext.Scraper().SetExtensionContext(&scraper.Context{
		Progress: func(message string) {
			m.status = message
		},
		Error: func(err error) {
			m.errorChan <- err
		},
	})
}

func (m *model) resetSpinner() {
	m.component.spinner = spinner.New()
	m.component.spinner.Spinner = spinner.Dot
}
