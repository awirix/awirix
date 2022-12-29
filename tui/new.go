package tui

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/vivi-app/vivi/stack"
	"golang.org/x/term"
	"os"
)

func newModel(options *Options) *model {
	if options == nil {
		options = &Options{}
	}

	m := &model{
		history:       stack.New[state](),
		selectedMedia: make(map[*lItem]struct{}),
		options:       options,
		errorChan:     make(chan error),
	}
	m.keyMap = NewKeyMap(m)

	m.current.state = stateExtensionSelect
	m.styles = DefaultStyles()

	newTextInput := func(placeholder string) textinput.Model {
		t := textinput.New()
		t.CharLimit = 80
		t.Placeholder = placeholder
		t.SetCursorMode(textinput.CursorStatic)
		return t
	}

	m.component.extensionSelect = newList("Extensions", "extension", "extensions")
	m.component.searchResults = newList("Search Results", "media", "media")
	m.component.textInput = newTextInput("Search...")
	m.component.actionSelect = newList("Actions", "action", "actions")
	m.component.mediaInfo = viewport.New(0, 0)
	m.component.help = help.New()

	m.text.mediaInfoTitle = "Info"

	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		width, height = 80, 24
	}
	m.resize(width, height)

	return m
}
