package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"golang.org/x/term"
	"os"
)

type model struct {
	options *Options
	error   chan error

	extensions []*extension.Extension
	current    struct {
		state     state
		extension *extension.Extension
		media     *scraper.Media
		layer     *scraper.Layer
		query     string
	}

	component struct {
		extensionSelect list.Model
	}

	style struct {
		global lipgloss.Style
	}
}

func (m *model) resize(width, height int) {
	frameX, frameY := m.style.global.GetFrameSize()
	m.component.extensionSelect.SetSize(width-frameX, height-frameY)
}

func newModel(options *Options) *model {
	if options == nil {
		options = &Options{}
	}

	model := &model{}

	model.current.state = stateExtensionSelect
	model.options = options
	model.error = make(chan error)
	model.style.global = lipgloss.NewStyle()

	newList := func(singular, plural string) list.Model {
		l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
		l.SetStatusBarItemName(singular, plural)
		return l
	}

	model.component.extensionSelect = newList("extension", "extensions")

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}
	model.resize(width, height)

	return model
}
