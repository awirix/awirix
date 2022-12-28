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

	Reset,
	Select,
	SelectAll,
	Info,
	Confirm,
	GoBack,

	Reverse key.Binding
}

func NewKeyMap() *KeyMap {
	return &KeyMap{
		Quit:      bind("Quit", "q"),
		ForceQuit: bind("Force Quit", "ctrl+c", "ctrl+d"),
		Reset:     bind("Reset", "backspace"),
		Select:    bind("Select", " "),
		SelectAll: bind("Select All", "tab"),
		Info:      bind("Info", "i"),
		Confirm:   bind("Confirm", "enter"),
		GoBack:    bind("Go Back", "esc"),
		Reverse:   bind("Reverse", "r"),
	}
}
