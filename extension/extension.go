package extension

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"os"
	"path/filepath"
	"strings"
)

type Extension struct {
	passport *passport.Passport
	scraper  *scraper.Scraper
}

func ListInstalled() []*Extension {
	extensions := make([]*Extension, 0)
	installed, err := filesystem.Api().ReadDir(where.Extensions())
	if err != nil {
		return extensions
	}

	for _, file := range installed {
		extension, err := FromPath(filepath.Join(where.Extensions(), file.Name()))
		if err != nil {
			continue
		}

		extensions = append(extensions, extension)
	}

	return extensions
}

func (e *Extension) LoadScraper() error {
	theScraper, err := scraper.FromPath(e.Path())
	if err != nil {
		return err
	}

	e.scraper = theScraper
	return nil
}

func FromPath(path string) (*Extension, error) {
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

	return New(thePassport, nil), nil
}

func New(passport *passport.Passport, scraper *scraper.Scraper) *Extension {
	return &Extension{
		passport: passport,
		scraper:  scraper,
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

	scraperTemplate, err := scraper.GenerateTemplate(domain)
	if err != nil {
		return nil, err
	}

	// TODO: dry
	err = filesystem.Api().WriteFile(filepath.Join(path, constant.Scraper), []byte(scraperTemplate), os.ModePerm)
	if err != nil {
		return nil, err
	}

	err = filesystem.Api().WriteFile(filepath.Join(path, constant.Test), []byte(scraper.TemplateTest), os.ModePerm)
	if err != nil {
		return nil, err
	}

	return FromPath(path)
}

func (e *Extension) String() string {
	return e.Passport().Name
}

func (e *Extension) Passport() *passport.Passport {
	return e.passport
}

func (e *Extension) Scraper() *scraper.Scraper {
	return e.scraper
}

func (e *Extension) IsInstalled() bool {
	path := e.Path()
	installed, err := filesystem.Api().ReadDir(where.Extensions())
	if err != nil {
		return false
	}

	for _, file := range installed {
		if filepath.Join(where.Extensions(), file.Name()) == path {
			return true
		}
	}

	return false
}

func (e *Extension) Path() string {
	dir := util.SanitizeFilename(e.Passport().ID)
	return filepath.Join(where.Extensions(), dir)
}

func (e *Extension) Install() error {
	svn, err := e.Passport().Github.Repository.SVNURL()
	if err != nil {
		return err
	}

	fmt.Println(svn)

	return nil
}

func (e *Extension) Uninstall() error {
	return filesystem.Api().RemoveAll(e.Path())
}
