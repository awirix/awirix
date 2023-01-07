package extension

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/awirix/awirix/extensions/passport"
	"github.com/awirix/awirix/filename"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/language"
	"github.com/awirix/awirix/where"
	"github.com/awirix/templates"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/samber/lo"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

func New(path string) (*Extension, error) {
	if !filepath.IsAbs(path) {
		var err error
		path, err = filepath.Abs(path)
		if err != nil {
			return nil, err
		}
	}

	ext := &Extension{
		path: path,
	}

	// extensions must have valid passports
	err := ext.loadPassport()
	if err != nil {
		return nil, err
	}

	return ext, nil
}

func (e *Extension) SetContext(ctx context.Context) {
	e.state.SetContext(ctx)
}

func GenerateInteractive() (*Extension, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, err
	}

	username := usr.Username

	if err != nil {
		return nil, err
	}

	answers := struct {
		Preset   string
		Name     string
		About    string
		Nsfw     bool
		Tags     []string
		Language string
	}{}

	err = survey.Ask([]*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Extension name",
			},
		},
		{
			Name: "about",
			Prompt: &survey.Input{
				Message: "Extension description",
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
				Message: "Extension tags (comma separated)",
			},
			Validate: func(val any) error {
				v := strings.Trim(val.(string), ", ")

				if len(strings.Split(v, ",")) > 5 {
					return errors.New("extension cannot have more than 5 tags")
				}

				return nil
			},
			Transform: func(ans any) (newAns any) {
				var tags []string

				for _, tag := range strings.Split(strings.Trim(ans.(string), ", "), ",") {
					if tag == "" {
						continue
					}

					tags = append(tags, strings.TrimSpace(tag))
				}

				return tags
			},
		},
		{
			Name: "language",
			Prompt: &survey.Select{
				Message:  "Extension language",
				Options:  language.Names,
				PageSize: 10,
				Default:  "English",
				Filter: func(filter string, value string, index int) bool {
					if fuzzy.MatchFold(filter, value) {
						return true
					}

					return fuzzy.MatchFold(filter, language.NativeNames[index])
				},
			},
			Transform: func(ans any) (newAns any) {
				name := ans.(survey.OptionAnswer).Value
				lang, _ := language.FromName(name)
				return survey.OptionAnswer{Value: lang.Code}
			},
		},
		{
			Name: "Preset",
			Prompt: &survey.Select{
				Message: "Programming language preset",
				Options: lo.Map(templates.PresetValues(), func(p templates.Preset, _ int) string {
					return p.String()
				}),
				Default: 0,
				VimMode: true,
			},
		},
	}, &answers)

	if err != nil {
		return nil, err
	}

	p := &passport.Passport{
		Name:        answers.Name,
		ID:          passport.ToID(answers.Name),
		About:       answers.About,
		VersionRaw:  "0.1.0",
		Tags:        answers.Tags,
		LanguageRaw: answers.Language,
		NSFW:        answers.Nsfw,
	}

	path := filepath.Join(where.Extensions(), username, filename.Sanitize(p.ID))

	exists, err := filesystem.Api().Exists(path)
	if err != nil {
		return nil, err
	}

	if exists {
		var overwrite bool
		err = survey.AskOne(&survey.Confirm{
			Message: "Extension with the same name already exists, overwrite?",
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

	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(p)
	if err != nil {
		return nil, err
	}

	preset, err := templates.PresetString(answers.Preset)
	if err != nil {
		return nil, err
	}

	tree, err := templates.Get(preset, templates.Info{
		Name:  p.Name,
		About: p.About,
		NSFW:  p.NSFW,
	})

	if err != nil {
		return nil, err
	}

	tree[filename.Passport] = &buffer

	for name, contents := range tree {
		err = filesystem.Api().WriteFile(filepath.Join(path, name), contents.Bytes(), os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	return New(path)
}
