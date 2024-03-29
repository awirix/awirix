package mini

import (
	"fmt"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/core"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/icon"
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
	LastSelectedSearchMedia *core.Media
	Action string
}

func stateSelectExtension(s *state) (err error) {
	exts, err := manager.Installed()
	if err != nil {
		return err
	}

	s.Extension, err = selectOne[*extension.Extension]("Select an extension to use", exts, renderExtension)
	if err != nil {
		return err
	}

	progress("Loading core")
	err = s.Extension.LoadCore(s.Options.Debug)
	if err != nil {
		return err
	}

	s.Extension.Core().SetExtensionContext(&core.Context{
		Progress: progress,
		Error: func(err error) {
			// TODO: handle it better
			progress(err.Error())
		},
	})

	if s.Extension.Core().HasSearch() {
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
	search := s.Extension.Core().Search()
	medias, err := search.Call(s.Query)
	if err != nil {
		return err
	}

	if len(medias) == 0 {
		notFound()
		return stateInputQuery(s)
	}

	s.LastSelectedMedia, err = selectOne[*core.Media](search.String(), medias, renderMedia)
	if err != nil {
		return err
	}
	s.LastSelectedSearchMedia = s.LastSelectedMedia

	if s.Extension.Core().HasLayers() {
		return stateLayers(s)
	}

	return stateDoAction(s)
}

func stateLayers(s *state) error {
	layers := s.Extension.Core().Layers()

	for _, layer := range layers {
		medias, err := layer.Call(s.LastSelectedMedia)
		if err != nil {
			return err
		}

		if len(medias) == 0 {
			notFound()
			if s.Extension.Core().HasSearch() {
				return stateSearchMedia(s)
			}

			return fmt.Errorf("nothing was found")
		}

		s.LastSelectedMedia, err = selectOne[*core.Media](layer.String(), medias, renderMedia)
		if err != nil {
			return err
		}
	}

	return stateDoAction(s)
}

func stateDoAction(s *state) error {
	var actions = make([]string, 0)

	for _, action := range s.Extension.Core().Actions() {
		actions = append(actions, action.String())
	}

	action, err := selectOne("What do you want to do?", s.Extension.Core().Actions(), func(s *core.Action) string { return s.String() })
	if err != nil {
		return err
	}

	err = action.Call(s.LastSelectedMedia)
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
	if s.Extension.Core().HasSearch() {
		options = append(options, optionSearch)
	}

	var optionLayer string

	if s.Extension.Core().HasLayers() {
		layers := s.Extension.Core().Layers()
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
