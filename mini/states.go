package mini

import (
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/scraper"
	"github.com/awirix/awirix/style"
)

func notFound() {
	theSpinner.Stop()
	defer theSpinner.Start()
	fmt.Printf(
		"%s %s\n",
		style.New().Foreground(color.Red).Bold(true).Render(icon.Cross),
		style.Fg(color.Red)("not found, try again"),
	)
}

type state struct {
	Options   *Options
	Extension *extension.Extension
	Query     string
	LastSelectedMedia,
	LastSelectedSearchMedia *scraper.Media
	Action string
}

func stateSelectExtension(s *state) (err error) {
	exts, err := manager.InstalledExtensions()
	if err != nil {
		return err
	}

	s.Extension, err = selectOne[*extension.Extension]("Select an extension to use", exts, renderExtension)
	if err != nil {
		return err
	}

	progress("Loading scraper")
	err = s.Extension.LoadScraper(s.Options.Debug)
	if err != nil {
		return err
	}

	s.Extension.Scraper().SetExtensionContext(&scraper.Context{
		Progress: progress,
		Error: func(err error) {
			// TODO: handle it better
			progress(err.Error())
		},
	})

	if s.Extension.Scraper().HasSearch() {
		return stateInputQuery(s)
	}

	return stateLayers(s)
}

func stateInputQuery(s *state) (err error) {
	s.Query, err = getString("Enter a search query", "")
	if err != nil {
		return
	}

	return stateSearchMedia(s)
}

func stateSearchMedia(s *state) error {
	search := s.Extension.Scraper().Search()
	medias, err := search.Call(s.Query)
	if err != nil {
		return err
	}

	if len(medias) == 0 {
		notFound()
		return stateInputQuery(s)
	}

	s.LastSelectedMedia, err = selectOne[*scraper.Media](search.String(), medias, renderMedia)
	if err != nil {
		return err
	}
	s.LastSelectedSearchMedia = s.LastSelectedMedia

	if s.Extension.Scraper().HasLayers() {
		return stateLayers(s)
	}

	return stateDoAction(s)
}

func stateLayers(s *state) error {
	layers := s.Extension.Scraper().Layers()

	for _, layer := range layers {
		medias, err := layer.Call(s.LastSelectedMedia)
		if err != nil {
			return err
		}

		if len(medias) == 0 {
			notFound()
			if s.Extension.Scraper().HasSearch() {
				return stateSearchMedia(s)
			}

			return fmt.Errorf("nothing was found")
		}

		s.LastSelectedMedia, err = selectOne[*scraper.Media](layer.String(), medias, renderMedia)
		if err != nil {
			return err
		}
	}

	return stateDoAction(s)
}

func stateDoAction(s *state) error {
	var actions = make([]string, 0)

	for _, action := range s.Extension.Scraper().Actions() {
		actions = append(actions, action.String())
	}

	action, err := selectOne[*scraper.Action]("What do you want to do?", s.Extension.Scraper().Actions(), func(s *scraper.Action) string { return s.String() })
	if err != nil {
		return err
	}

	err = action.Call([]*scraper.Media{s.LastSelectedMedia})
	if err != nil {
		return err
	}

	return stateDoNext(s)
}

func stateDoNext(s *state) error {
	const (
		optionQuit            = "Quit"
		optionSelectExtension = "Select Extension"
		optionSearch          = "Search"
	)

	options := []string{optionSelectExtension}
	if s.Extension.Scraper().HasSearch() {
		options = append(options, optionSearch)
	}

	var optionLayer string

	if s.Extension.Scraper().HasLayers() {
		layers := s.Extension.Scraper().Layers()
		optionLayer = fmt.Sprintf(`Back to the "%s"`, layers[0].String())
		options = append(options, optionLayer)
	}

	options = append(options, optionQuit)

	clearScreen()
	option, err := selectOne[string]("Done! What to do next?", options, func(s string) string { return s })
	if err != nil {
		return err
	}

	switch option {
	case optionQuit:
		return nil
	case optionSelectExtension:
		return stateSelectExtension(s)
	case optionSearch:
		return stateInputQuery(s)
	case optionLayer:
		s.LastSelectedMedia = s.LastSelectedSearchMedia
		return stateLayers(s)
	default:
		return fmt.Errorf("unknown option %s", option)
	}
}
