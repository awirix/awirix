package tui

import (
	"fmt"
	zone "github.com/lrstanley/bubblezone"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/style"
)

func (m *model) View() string {
	switch m.current.state {
	case stateLoading:
		return m.status
	case stateError:
		return style.Fg(color.Red)(m.current.error.Error())
	case stateExtensionSelect:
		return zone.Scan(m.style.global.Render(m.component.extensionSelect.View()))
	case stateSearch:
		return m.component.textInput.View()
	default:
		panic(fmt.Sprintf(`Unknown state "%s"`, m.current.state.String()))
	}
}
