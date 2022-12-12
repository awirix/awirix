package extension

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/manifoldco/promptui"
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/filename"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/template"
	"github.com/vivi-app/vivi/where"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func New(path string) (*Extension, error) {
	ext := &Extension{path: path}
	err := ext.loadPassport()
	if err != nil {
		return nil, err
	}

	return ext, nil
}

func GenerateInteractive() (*Extension, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	username := usr.Username

	var (
		prompt promptui.Prompt
		sel    promptui.Select
	)

	answers := struct {
		Name     string
		About    string
		Nsfw     bool
		Tags     string
		Language string
	}{}

	prompt = promptui.Prompt{
		Label: "Name of the extension",
	}

	answers.Name, err = prompt.Run()
	if err != nil {
		return nil, err
	}

	prompt.Label = "About the extension"
	answers.About, err = prompt.Run()
	if err != nil {
		return nil, err
	}

	prompt.Label = "NSFW (18+)"
	prompt.IsConfirm = true
	prompt.Default = "N"

	_, err = prompt.Run()
	answers.Nsfw = err == nil

	prompt.Label = "Tags (separated by commas)"
	prompt.IsConfirm = false

	answers.Tags, err = prompt.Run()
	if err != nil {
		return nil, err
	}

	sel = promptui.Select{
		Label: "Language",
		Items: language.Names,
		Size:  7,
		Searcher: func(input string, index int) bool {
			// TODO: fuzzy search
			name := strings.ReplaceAll(strings.ToLower(language.Names[index]), " ", "")
			input = strings.ReplaceAll(strings.ToLower(input), " ", "")

			return strings.Contains(name, input)
		},
		CursorPos:         lo.IndexOf(language.Names, "English"),
		StartInSearchMode: true,
	}

	_, answers.Language, err = sel.Run()
	if err != nil {
		return nil, err
	}

	lang, _ := language.FromName(answers.Language)

	p := &passport.Passport{
		Name:        answers.Name,
		ID:          passport.ToID(answers.Name),
		About:       answers.About,
		VersionRaw:  "0.1.0",
		LanguageRaw: lang.Code,
		NSFW:        answers.Nsfw,
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

	path := filepath.Join(where.Extensions(), username, filename.Sanitize(p.ID))

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return nil, err
	}

	if exists {
		var overwrite bool
		err = survey.AskOne(&survey.Confirm{
			Message: "Extension already exists, overwrite?",
			Default: false,
		}, &overwrite)
		if err != nil {
			return nil, err
		}

		if !overwrite {
			return nil, fmt.Errorf("cancelled")
		}

		err = filesystem.Api().RemoveAll(path)
		if err != nil {
			return nil, err
		}
	}

	err = filesystem.Api().MkdirAll(path, os.ModePerm)
	if err != nil {
		return nil, err
	}

	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(p)
	if err != nil {
		return nil, err
	}

	for _, t := range []struct {
		filename string
		contents []byte
	}{
		{filename.Passport, data.Bytes()},
		{filename.Scraper, template.Scraper()},
		{filename.Tester, template.Tester()},
		{filename.EditorConfig, template.EditorConfig()},
	} {
		err = filesystem.Api().WriteFile(filepath.Join(path, t.filename), t.contents, os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return New(path)
}
