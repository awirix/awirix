package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/context"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/tui/bind"
	"golang.org/x/term"
	"os"
)

func newModel(options *Options) *model {
	if options == nil {
		options = &Options{}
	}

	model := &model{
		keyMap:  bind.NewKeyMap(),
		history: stack.New[state](),
		context: context.New(),
	}

	model.current.state = stateExtensionSelect
	model.options = options
	model.error = make(chan error)
	model.style.global = lipgloss.NewStyle()

	newList := func(singular, plural string) list.Model {
		l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
		l.SetStatusBarItemName(singular, plural)
		return l
	}

	newTextInput := func(placeholder string) textinput.Model {
		t := textinput.New()
		t.CharLimit = 80
		t.Placeholder = placeholder
		return t
	}

	model.component.extensionSelect = newList("extension", "extensions")
	model.component.layerResults = newList("media", "medias")
	model.component.searchResults = newList("media", "medias")
	model.component.textInput = newTextInput("Search...")

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}
	model.resize(width, height)

	return model
}
