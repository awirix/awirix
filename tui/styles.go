package tui

import (
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/style"
)

type Styles struct {
	global,
	title,
	titleError,
	titleBar lipgloss.Style
}

func DefaultStyles() (s Styles) {
	listStyles := list.DefaultStyles()
	s.global = style.New()
	s.title = listStyles.Title
	s.titleError = s.title.Copy().Background(color.Red)
	s.titleBar = listStyles.TitleBar

	return
}
