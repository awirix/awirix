package mini

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/extensions/manager"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/style"
	"golang.org/x/exp/slices"
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

	s.Extension.Scraper().SetProgress(progress)

	if s.Options.EditConfig {
		return stateExtensionConfig(s)
	}

	if s.Extension.Scraper().HasSearch() {
		return stateInputQuery(s)
	}

	return stateLayers(s)
}

func stateExtensionConfig(s *state) (err error) {
	type sectionWithName *lo.Tuple2[string, *passport.ConfigSection]
	var sections = make([]sectionWithName, 0)
	for name, section := range s.Extension.Passport().Config {
		sections = append(sections, &lo.Tuple2[string, *passport.ConfigSection]{A: name, B: section})
	}

	slices.SortFunc(sections, func(a, b sectionWithName) bool {
		return a.A < b.A
	})

	for _, t := range sections {
		section := t.B
		value, err := getConfigValue(section)
		if err != nil {
			return err
		}

		section.SetValue(value)
	}

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

	s.LastSelectedMedia, err = selectOne[*scraper.Media](search.Title(), medias, renderMedia)
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

		s.LastSelectedMedia, err = selectOne[*scraper.Media](layer.Title(), medias, renderMedia)
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

	prepared, err := s.Extension.Scraper().Prepare(s.LastSelectedMedia)
	if err != nil {
		return err
	}
	s.LastSelectedMedia = prepared

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

	return stateDoNext(s)
}

func stateDoNext(s *state) error {
	const (
		optionQuit            = "Quit"
		optionSelectExtension = "Select Extension"
		optionSearch          = "Search"
		optionStream          = "Stream"
		optionDownload        = "Download"
	)

	options := []string{optionSelectExtension}
	if s.Extension.Scraper().HasSearch() {
		options = append(options, optionSearch)
	}

	var optionLayer string

	if s.Extension.Scraper().HasLayers() {
		layers := s.Extension.Scraper().Layers()
		optionLayer = fmt.Sprintf(`Back to the "%s"`, layers[0].Title())
		options = append(options, optionLayer)
	}

	options = append(options, optionQuit, optionStream, optionDownload)

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
	case optionStream:
		err = s.Extension.Scraper().Stream(s.LastSelectedMedia)
		if err != nil {
			return err
		}

		return stateDoNext(s)
	case optionDownload:
		err = s.Extension.Scraper().Download(s.LastSelectedMedia)
		if err != nil {
			return err
		}

		return stateDoNext(s)
	default:
		return fmt.Errorf("unknown option %s", option)
	}
}
