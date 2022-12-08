package executil

import (
	"os/exec"
	"strings"
)

func ProgramInPath(program string) bool {
	_, err := exec.LookPath(strings.TrimSpace(program))
	return err == nil
}
