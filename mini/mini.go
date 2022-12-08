package mini

import (
	"fmt"
	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/vivi-app/vivi/extensions/manager"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
	"os"
	"strings"
	"time"
)

var theSpinner = spinner.New(spinner.CharSets[9], 100*time.Millisecond, spinner.WithWriter(os.Stderr), spinner.WithColor("cyan"))

func Run(options *Options) error {
	progress := func(text string) {
		theSpinner.Suffix = " " + text
	}

	progress("Searching for installed extensions...")
	theSpinner.Start()
	defer theSpinner.Stop()

	if options == nil {
		options = &Options{}
	}

	extensions, err := manager.InstalledExtensions()
	if err != nil {
		return err
	}

	theSpinner.Stop()

	promptSelect := promptui.Select{
		Label:     "Select an extension to use",
		Items:     extensions,
		Size:      7,
		IsVimMode: true,
		Stdout:    os.Stderr,
	}

	index, _, err := promptSelect.Run()
	if err != nil {
		return err
	}

	progress("Loading extension...")
	theSpinner.Start()

	ext := extensions[index]

	ext.Init()

	progress("Loading extension's data...")

	if err = ext.LoadPassport(); err != nil {
		return err
	}

	if err = ext.LoadScraper(); err != nil {
		return err
	}

	theScraper := ext.Scraper()
	theScraper.SetProgress(progress)

	theSpinner.Stop()
	promptInput := promptui.Prompt{
		Label:       "Enter a search query",
		IsVimMode:   true,
		HideEntered: true,
		Stdout:      os.Stderr,
	}

	query, err := promptInput.Run()
	if err != nil {
		return err
	}
	theSpinner.Start()

	medias, err := theScraper.Search(query)
	if err != nil {
		return err
	}

	theSpinner.Stop()
	if len(medias) == 0 {
		return fmt.Errorf("no results found")
	}
	promptSelect = promptui.Select{
		Label:     "Select a media",
		Items:     mediasToStringSlice(medias),
		Size:      7,
		IsVimMode: true,
		Stdout:    os.Stderr,
	}

	index, _, err = promptSelect.Run()
	if err != nil {
		return err
	}

	theSpinner.Start()
	media := medias[index]

	medias, err = theScraper.Explore(media)
	if err != nil {
		return err
	}

	theSpinner.Stop()
	if len(medias) == 0 {
		return fmt.Errorf("no results found")
	}

	promptSelect = promptui.Select{
		Label:     "Select a media",
		Items:     mediasToStringSlice(medias),
		Size:      7,
		IsVimMode: true,
		Stdout:    os.Stderr,
	}

	index, _, err = promptSelect.Run()
	if err != nil {
		return err
	}

	theSpinner.Start()
	media = medias[index]

	media, err = theScraper.Prepare(media)
	if err != nil {
		return err
	}

	theSpinner.Stop()
	const (
		actionStream   = "Stream"
		actionDownload = "Download"
	)

	promptSelect = promptui.Select{
		Label:  "What do you want to do?",
		Items:  []string{actionStream, actionDownload},
		Stdout: os.Stderr,
	}

	_, action, err := promptSelect.Run()
	if err != nil {
		return err
	}

	switch action {
	case actionStream:
		return theScraper.Stream(media)
	case actionDownload:
		return theScraper.Download(media)
	default:
		return fmt.Errorf("unknown action: %s", action)
	}
}

func mediasToStringSlice(medias []*scraper.Media) []string {
	slice := make([]string, len(medias))

	for i, media := range medias {
		var b strings.Builder
		b.WriteString(media.String())
		if media.HasAbout() {
			b.WriteString(" ")
			b.WriteString(style.Faint(media.About()))
		}

		slice[i] = b.String()
	}

	return slice
}
