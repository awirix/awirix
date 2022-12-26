package bind

import (
	"github.com/charmbracelet/bubbles/key"
)

func bind(help string, keys ...string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], help),
	)
}

type KeyMap struct {
	Quit, ForceQuit,

	Select, Confirm, GoBack,

	Reverse key.Binding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		Quit:      bind("Quit", "q"),
		ForceQuit: bind("Force Quit", "ctrl+c", "ctrl+d"),
		Select:    bind("Select", " "),
		Confirm:   bind("Confirm", "enter"),
		GoBack:    bind("Go Back", "esc"),
		Reverse:   bind("Reverse", "r"),
	}
}
