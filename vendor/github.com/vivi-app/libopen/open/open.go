package open

// Run open a file, directory, or URI using the OS's default application for that object type.
// Wait for the open command to complete.
func Run(input string) error {
	return open(input).Run()
}

// Start Open a file, directory, or URI using the OS's default application for that object type.
// Don't wait for the open command to complete.
func Start(input string) error {
	return open(input).Start()
}

// RunWith Open a file, directory, or URI using the specified application.
// Wait for the open command to complete.
func RunWith(input, appName string) error {
	return openWith(input, appName).Run()
}

// StartWith open a file, directory, or URI using the specified application.
// Don't wait for the open command to complete.
func StartWith(input, appName string) error {
	return openWith(input, appName).Start()
}
