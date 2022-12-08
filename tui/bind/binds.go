package bind

import "github.com/charmbracelet/bubbles/key"

func bind(help string, keys ...string) key.Binding {
	return key.NewBinding(
		key.WithKeys(keys...),
		key.WithHelp(keys[0], help),
	)
}

var (
	Confirm   = bind("Confirm", "enter")
	Cancel    = bind("Cancel", "esc")
	Quit      = bind("Quit", "q")
	ForceQuit = bind("Force Quit", "ctrl+c", "ctrl+d")
)
