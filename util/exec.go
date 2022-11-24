package util

import (
	"os/exec"
	"strings"
)

func ProgramInPath(program string) bool {
	if _, err := exec.LookPath(strings.TrimSpace(program)); err != nil {
		return false
	}

	return true
}
