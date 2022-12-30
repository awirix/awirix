package tui

import (
	"context"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/style"
	"golang.org/x/exp/slices"
	"strings"
)

type model struct {
	options   *Options
	errorChan chan error
	history   *stack.Stack[state]

	extensions []*extension.Extension

	// algebraic set
	selectedMedia map[*lItem]struct{}

	current struct {
		dimensions struct {
			terminalWidth, terminalHeight   int
			availableWidth, availableHeight int
		}
		state         state
		extension     *extension.Extension
		layer         *scraper.Layer
		error         error
		context       context.Context
		cancelContext context.CancelFunc
		mediaInfo     string
	}

	component struct {
		extensionSelect list.Model
		textInput       textinput.Model
		searchResults   list.Model
		layers          map[string]*list.Model
		actionSelect    list.Model
		spinner         spinner.Model
		mediaInfo       viewport.Model
		help            help.Model
	}

	text struct {
		searchTitle    string
		mediaInfoTitle string
		mediaInfoName  string
		status         string
	}

	keyMap *KeyMap

	styles Styles
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
	frameX, frameY := m.styles.global.GetFrameSize()

	m.current.dimensions.terminalWidth, m.current.dimensions.terminalHeight = width, height
	m.current.dimensions.availableWidth, m.current.dimensions.availableHeight = width-frameX, height-frameY

	for _, lst := range m.lists() {
		lst.SetSize(
			m.current.dimensions.availableWidth,
			m.current.dimensions.availableHeight,
		)
	}

	mediaInfoHeaderHeight := lipgloss.Height(m.styles.titleBar.Render(m.styles.title.Render(m.text.mediaInfoTitle))) + lipgloss.Height(m.styles.statusBar.Render(m.text.mediaInfoName))
	helpHeight := lipgloss.Height(m.styles.helpStyle.Render(m.component.help.View(m.keyMap)))

	m.component.mediaInfo.Height = m.current.dimensions.availableHeight - mediaInfoHeaderHeight - helpHeight
	m.component.mediaInfo.Width = m.current.dimensions.availableWidth

	// error can not occur here
	r, _ := glamour.NewTermRenderer(glamour.WithAutoStyle(), glamour.WithWordWrap(m.component.mediaInfo.Width))
	// but it can here
	md, err := r.Render(m.current.mediaInfo)
	if err != nil {
		md = m.current.mediaInfo
	}

	m.component.mediaInfo.SetContent(strings.TrimSpace(md))
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
	ctx := context.WithValue(context.Background(), "extension", m.current.extension)
	m.current.context, m.current.cancelContext = context.WithCancel(ctx)
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
			m.text.status = message
		},
		Error: func(err error) {
			m.errorChan <- err
		},
	})
}

func (m *model) resetSpinner() {
	m.component.spinner = spinner.New()
	m.component.spinner.Style = style.New().Foreground(color.Purple)
	m.component.spinner.Spinner = spinner.Dot
}
