package mini

import (
	"fmt"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/extensions/manager"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
)

func notFound() {
	fmt.Printf("%s not found, try again\n", style.Fg(color.Red)(icon.Cross))
}

type state struct {
	Extension *extension.Extension
	Query     string
	LastSelectedMedia,
	LastSelectedSearchMedia *scraper.Media
	Action string
}

func stateSelectExtension() (err error) {
	exts, err := manager.InstalledExtensions()
	if err != nil {
		return err
	}

	s := &state{}

	s.Extension, err = selectOne[*extension.Extension]("Select an extension to use", exts, renderExtension)
	if err != nil {
		return err
	}

	err = s.Extension.LoadScraper()
	if err != nil {
		return err
	}

	if s.Extension.Scraper().HasSearch() {
		return stateInputQuery(s)
	}

	return stateLayers(s)
}

func stateInputQuery(s *state) (err error) {
	s.Query, err = getString("Enter a search query")
	if err != nil {
		return
	}

	return stateSearchMedia(s)
}

func stateSearchMedia(s *state) error {
	medias, err := s.Extension.Scraper().Search(s.Query)
	if err != nil {
		return err
	}

	if len(medias) == 0 {
		notFound()
		return stateInputQuery(s)
	}

	s.LastSelectedMedia, err = selectOne[*scraper.Media]("Select a media", medias, renderMedia)
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
	layers, err := s.Extension.Scraper().Layers()
	if err != nil {
		return err
	}

	for _, layer := range layers {
		medias, err := layer.Function(s.LastSelectedMedia)
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

		s.LastSelectedMedia, err = selectOne[*scraper.Media](layer.Name, medias, renderMedia)
		if err != nil {
			return err
		}
	}

	return stateDoAction(s)
}

func stateDoAction(s *state) error {
	const (
		actionStream   = "Stream"
		actionDownload = "Download"
	)

	var actions = make([]string, 0)

	if s.Extension.Scraper().HasStream() {
		actions = append(actions, actionStream)
	}

	if s.Extension.Scraper().HasDownload() {
		actions = append(actions, actionDownload)
	}

	action, err := selectOne[string]("What do you want to do?", actions, func(s string) string { return s })
	if err != nil {
		return err
	}

	switch action {
	case actionStream:
		err = s.Extension.Scraper().Stream(s.LastSelectedMedia)
		if err != nil {
			return err
		}
	case actionDownload:
		err = s.Extension.Scraper().Download(s.LastSelectedMedia)
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unknown action %s", action)
	}

	const (
		optionQuit            = "Quit"
		optionSelectExtension = "Select Extension"
		optionSearchNew       = "Search New"
		optionSearch          = "Search"
		optionLayer           = "Back to the first layer" // TODO: use actual name of the first layer
		// TODO: stream & download
	)

	options := []string{optionSelectExtension}
	if s.Extension.Scraper().HasSearch() {
		options = append(options, optionSearchNew, optionSearch)
	}

	if s.Extension.Scraper().HasLayers() {
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
		return stateSelectExtension()
	case optionSearchNew:
		return stateInputQuery(s)
	case optionSearch:
		return stateSearchMedia(s)
	case optionLayer:
		s.LastSelectedMedia = s.LastSelectedSearchMedia
		return stateLayers(s)
	default:
		return fmt.Errorf("unknown option %s", option)
	}
}
