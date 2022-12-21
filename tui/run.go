package tui

import (
	tea "github.com/charmbracelet/bubbletea"
	zone "github.com/lrstanley/bubblezone"
)

func Run(options *Options) error {
	m := newModel(options)

	var program *tea.Program

	zone.NewGlobal()
	defer zone.Close()

	if options.AltScreen {
		program = tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	} else {
		program = tea.NewProgram(m, tea.WithMouseCellMotion())
	}

	_, err := program.Run()
	return err
}
