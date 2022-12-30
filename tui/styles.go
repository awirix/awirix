package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/viper"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/key"
	"github.com/vivi-app/vivi/style"
)

type Styles struct {
	global,
	helpStyle,
	title,
	titleError,
	titleBar,
	statusBar lipgloss.Style
}

func DefaultStyles() (s Styles) {
	listStyles := list.DefaultStyles()

	s.global = style.New().Padding(
		viper.GetInt(key.TUIPaddingTop),
		viper.GetInt(key.TUIPaddingRight),
		viper.GetInt(key.TUIPaddingBottom),
		viper.GetInt(key.TUIPaddingLeft),
	)

	s.title = listStyles.Title
	s.titleError = s.title.Copy().Background(color.Red)
	s.titleBar = listStyles.TitleBar
	s.statusBar = listStyles.StatusBar
	s.helpStyle = listStyles.HelpStyle

	return
}
