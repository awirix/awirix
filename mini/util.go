package mini

import (
	"os"
	"os/exec"
	"runtime"
)

func clearScreen() {
	run := func(name string, args ...string) error {
		command := exec.Command(name, args...)
		command.Stdout = os.Stdout
		command.Stdin = os.Stdin
		command.Stderr = os.Stderr
		return command.Run()
	}

	switch runtime.GOOS {
	case "linux", "darwin":
		err := run("tput", "clear")
		if err != nil {
			_ = run("clear")
		}
	case "windows":
		_ = run("cls")
	}
}
