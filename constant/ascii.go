package constant

import (
	_ "embed"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/style"
	"strings"
)

//go:embed spider-ascii-art.txt
var AsciiArt string

func init() {
	// make spider eyes red
	AsciiArt = strings.ReplaceAll(AsciiArt, `"`, style.Fg(color.Red)(`"`))
}
