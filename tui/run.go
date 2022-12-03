package tui

import tea "github.com/charmbracelet/bubbletea"

func Run(options *Options) error {
	if options == nil {
		options = &Options{}
	}

	m := &model{options: options}

	var program *tea.Program

	if options.AltScreen {
		program = tea.NewProgram(m, tea.WithAltScreen())
	} else {
		program = tea.NewProgram(m)
	}

	_, err := program.Run()
	return err
}
