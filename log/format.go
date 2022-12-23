package log

import (
	"fmt"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/style"
	"io"
)

func WriteErrorf(out io.Writer, format string, args ...any) (int, error) {
	a, err := out.Write([]byte(style.New().Foreground(color.Red).Bold(true).Render(icon.Cross)))
	if err != nil {
		return a, err
	}

	b, err := out.Write([]byte{0x20})
	if err != nil {
		return a + b, err
	}

	c, err := out.Write([]byte(style.Fg(color.Red)(fmt.Sprintf(format, args...))))
	return a + b + c, err
}

func WriteSuccessf(out io.Writer, format string, args ...any) (int, error) {
	a, err := out.Write([]byte(style.New().Foreground(color.Green).Bold(true).Render(icon.Check)))
	if err != nil {
		return a, err
	}

	b, err := out.Write([]byte{0x20})
	if err != nil {
		return a + b, err
	}

	c, err := out.Write([]byte(style.Fg(color.Green)(fmt.Sprintf(format, args...))))
	return a + b + c, err
}
