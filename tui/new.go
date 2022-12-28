package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/stack"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/tui/bind"
	"golang.org/x/term"
	"os"
)

func newModel(options *Options) *model {
	if options == nil {
		options = &Options{}
	}

	model := &model{
		keyMap:        bind.NewKeyMap(),
		history:       stack.New[state](),
		selectedMedia: make(map[*lItem]struct{}),
		options:       options,
		errorChan:     make(chan error),
	}

	listStyles := list.DefaultStyles()
	model.current.state = stateExtensionSelect
	model.style.global = style.New()
	model.style.title = listStyles.Title
	model.style.titleError = model.style.title.Copy().Background(color.Red)
	model.style.titleBar = listStyles.TitleBar

	newTextInput := func(placeholder string) textinput.Model {
		t := textinput.New()
		t.CharLimit = 80
		t.Placeholder = placeholder
		t.SetCursorMode(textinput.CursorStatic)
		return t
	}

	model.component.extensionSelect = newList("Extensions", "extension", "extensions")
	model.component.searchResults = newList("Search Results", "media", "media")
	model.component.textInput = newTextInput("Search...")
	model.component.actionSelect = newList("Actions", "action", "actions")

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}
	model.resize(width, height)

	return model
}
