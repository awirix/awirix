//go:build darwin
// +build darwin

package open

import (
	"os/exec"
)

func open(input string) *exec.Cmd {
	return exec.Command("open", input)
}

func openWith(input, appName string) *exec.Cmd {
	return exec.Command("open", "-a", appName, input)
}
