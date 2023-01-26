package tui

import (
	"github.com/awirix/awirix/key"
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
	"github.com/spf13/viper"
)

func Run(options *Options) error {
	m := newModel(options)

	zone.NewGlobal()
	defer zone.Close()

	var programOptions = make([]tea.ProgramOption, 0)

	if viper.GetBool(key.TUIClickable) {
		programOptions = append(programOptions, tea.WithMouseCellMotion())
	}

	if options.AltScreen {
		programOptions = append(programOptions, tea.WithAltScreen())
	}

	_, err := tea.NewProgram(m, programOptions...).Run()
	return err
}
