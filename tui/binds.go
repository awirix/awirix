package tui

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
	model *model

	Quit, ForceQuit,

	Reset,
	Select,
	SelectAll,
	Info,
	Confirm,
	GoBack,

	Reverse key.Binding
}

func NewKeyMap(m *model) *KeyMap {
	return &KeyMap{
		model:     m,
		Quit:      bind("quit", "q"),
		ForceQuit: bind("force quit", "ctrl+c", "ctrl+d"),
		Reset:     bind("reset", "backspace"),
		Select:    bind("select", " "),
		SelectAll: bind("select all", "tab"),
		Info:      bind("info", "i"),
		Confirm:   bind("confirm", "enter"),
		GoBack:    bind("back", "esc"),
		Reverse:   bind("reverse", "r"),
	}
}
