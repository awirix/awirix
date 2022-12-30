package tui

import (
	"github.com/charmbracelet/bubbles/key"
)

func (k *KeyMap) ShortHelp() []key.Binding {
	l := func(binds ...key.Binding) []key.Binding {
		return append(binds, k.GoBack)
	}

	switch k.model.current.state {
	case stateLoading:
		return l()
	case stateError:
		return l(k.Quit)
	case stateExtensionSelect:
		return nil
	case stateSearch:
		return l(k.Confirm)
	case stateSearchResults, stateLayer, stateActionSelect:
		return l(k.Select, k.Confirm, k.Reverse)
	case stateMediaInfo:
		return l()
	default:
		return nil
	}
}

func (k *KeyMap) FullHelp() [][]key.Binding {
	//TODO: Add more help
	return [][]key.Binding{k.ShortHelp()}
}
