//go:build android
// +build android

package open

func open(input string) *exec.Cmd {
	return exec.Command("termux-open", input)
}

func openWith(input, _ string) *exec.Cmd {
	return exec.Command("termux-open", "--choose", input)
}
