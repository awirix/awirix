package extension

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/template"
	"github.com/vivi-app/vivi/tester"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"os"
	"path/filepath"
	"strings"
)

func NewFromID(id string) (*Extension, bool) {
	extensions := ListInstalled()
	for _, extension := range extensions {
		if extension.Passport().ID == id {
			return extension, true
		}
	}

	return nil, false
}

func NewFromPath(path string) (*Extension, error) {
	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return nil, err
	}

	if !exists {
		return nil, fmt.Errorf("path does not exist: %s", path)
	}

	isDir, err := filesystem.Api().IsDir(path)
	if err != nil {
		return nil, err
	}

	if !isDir {
		return nil, fmt.Errorf("path is not a directory: %s", path)
	}

	thePassport, err := passport.FromPath(path)
	if err != nil {
		return nil, err
	}

	return New(thePassport, nil, nil), nil
}

func New(passport *passport.Passport, scraper *scraper.Scraper, tester *tester.Tester) *Extension {
	return &Extension{
		passport: passport,
		scraper:  scraper,
		tester:   tester,
	}
}

func NewInteractive() (*Extension, error) {
	var questions = []*survey.Question{
		{
			Name: "domain",
			Prompt: &survey.Select{
				Message: "Domain",
				Options: []string{string(passport.DomainAnime), string(passport.DomainMovies), string(passport.DomainShows)},
				VimMode: true,
			},
			Validate: survey.Required,
		},
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Name of the extension",
			},
			Validate: survey.Required,
		},
		{
			Name: "about",
			Prompt: &survey.Input{
				Message: "About the extension",
			},
			Validate: survey.MaxLength(100),
		},
		{
			Name: "nsfw",
			Prompt: &survey.Confirm{
				Message: "Is this extension NSFW?",
			},
		},
		{
			Name: "tags",
			Prompt: &survey.Input{
				Message: "Tags (comma separated)",
			},
		},
		{
			Name: "language",
			Prompt: &survey.Select{
				Renderer: survey.Renderer{},
				Message:  "Language of the extension",
				Options:  language.Names,
				Default:  "English",
				PageSize: 8,
			},
		},
	}

	answers := struct {
		Domain   string
		Name     string
		About    string
		Nsfw     bool
		Tags     string
		Language string
	}{}

	err := survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	lang, _ := language.FromName(answers.Language)

	var domain = passport.Domain(answers.Domain)

	p := &passport.Passport{
		Name:     answers.Name,
		ID:       passport.ToID(answers.Name),
		About:    answers.About,
		Version:  semver.MustParse("0.1.0"),
		Domain:   domain,
		Language: lang,
		NSFW:     answers.Nsfw,
		Config: map[string]*passport.ConfigSection{
			"test": {
				Display: "this is a test",
				About:   "about the test",
				Default: false,
			},
		},
	}

	if answers.Tags != "" {
		for _, tag := range strings.Split(answers.Tags, ",") {
			p.Tags = append(p.Tags, strings.TrimSpace(tag))
		}
	}

	path := filepath.Join(where.Extensions(), util.SanitizeFilename(p.ID))
	err = filesystem.Api().MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	data, err := toml.Marshal(p)
	if err != nil {
		return nil, err
	}

	err = filesystem.Api().WriteFile(filepath.Join(path, constant.Passport), data, os.ModePerm)
	if err != nil {
		return nil, err
	}

	for _, t := range []lo.Tuple2[string, func() []byte]{
		{constant.Scraper, func() []byte { return template.NewScraper(domain) }},
		{constant.Tester, template.NewTest},
	} {
		err = filesystem.Api().WriteFile(filepath.Join(path, t.A), t.B(), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return NewFromPath(path)
}