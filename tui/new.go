package tui

import (
	"github.com/awirix/awirix/stack"
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/viewport"
	"golang.org/x/term"
	"os"
	"time"
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

	minute := time.Minute
	m.component.extensionSelect = m.newList("Extensions", "extension", "extensions", &minute)
	m.component.searchResults = m.newList("Search Results", "media", "media", nil)
	m.component.textInput = newTextInput("Search...")
	m.component.actionSelect = m.newList("️Actions", "action", "actions", nil)
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
