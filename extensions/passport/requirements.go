package passport

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/util"
	"runtime"
	"strings"
)

type Requirements struct {
	OS       []string `toml:"os,omitempty"`
	Programs []string `toml:"programs,omitempty"`
}

func (d *Requirements) Info() string {
	var b strings.Builder

	if len(d.OS) != 0 {
		b.WriteString("OS ")

		if lo.Contains(d.OS, runtime.GOOS) {
			b.WriteString(style.Fg(color.Green)(icon.Check))
		} else {
			b.WriteString(style.Fg(color.Red)(icon.Cross))
		}

		b.WriteString(style.Faint(" Available on: "))
		b.WriteString(style.Faint(strings.Join(d.OS, ", ")))
		b.WriteString(style.Faint(fmt.Sprintf(". You're on: %s", runtime.GOOS)))

		b.WriteRune('\n')
	}

	if len(d.Programs) != 0 {
		b.WriteString("Programs ")

		for _, program := range d.Programs {
			if util.ProgramInPath(program) {
				b.WriteString(style.Fg(color.Green)(icon.Check))
			} else {
				b.WriteString(style.Fg(color.Red)(icon.Cross))
			}

			b.WriteRune(' ')
			b.WriteString(style.Faint(program))
			b.WriteRune(' ')
		}
	}

	return b.String()
}

func (d *Requirements) Matches() bool {
	if len(d.OS) != 0 && !lo.Contains(d.OS, runtime.GOOS) {
		return false
	}

	for _, program := range d.Programs {
		if !util.ProgramInPath(program) {
			return false
		}
	}

	return true
}
