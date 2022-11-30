package extension

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/template"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func New(path string) *Extension {
	return &Extension{path: path}
}

func GenerateInteractive() (*Extension, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	username := usr.Username

	var questions = []*survey.Question{
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

	err = survey.Ask(questions, &answers)
	if err != nil {
		return nil, err
	}

	lang, _ := language.FromName(answers.Language)

	p := &passport.Passport{
		Name:     answers.Name,
		ID:       passport.ToID(answers.Name),
		About:    answers.About,
		Version:  semver.MustParse("0.1.0"),
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

	path := filepath.Join(where.Extensions(), username, util.SanitizeFilename(p.ID))
	err = filesystem.Api().MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	data, err := toml.Marshal(p)
	if err != nil {
		return nil, err
	}

	for _, t := range []lo.Tuple2[string, []byte]{
		{constant.FilenamePassport, data},
		{constant.FilenameScraper, template.NewScraper()},
		{constant.FilenameTester, template.NewTest()},
	} {
		err = filesystem.Api().WriteFile(filepath.Join(path, t.A), t.B, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return New(path), nil
}
