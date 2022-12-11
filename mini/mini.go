package mini

import (
	"github.com/briandowns/spinner"
	"os"
	"time"
)

var theSpinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr), spinner.WithColor("cyan"))

func progress(message string) {
	theSpinner.Suffix = " " + message
}

func Run(options *Options) error {
	progress("Searching for installed extensions...")
	theSpinner.Start()
	defer theSpinner.Stop()

	if options == nil {
		options = &Options{}
	}

	return stateSelectExtension()
}
