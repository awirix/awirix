package app

import (
	_ "embed"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/style"
	"strings"
)

//go:embed spider-ascii-art.txt
var AsciiArt string

func init() {
	// make spider eyes red
	AsciiArt = strings.ReplaceAll(AsciiArt, `"`, style.Fg(color.Red)(`"`))
}
