package cmd

import (
	"errors"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/language"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/scraper/anime"
	"github.com/vivi-app/vivi/scraper/movie"
	"github.com/vivi-app/vivi/semver"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/util"
	"github.com/vivi-app/vivi/where"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(extensionsCmd)
}

var extensionsCmd = &cobra.Command{
	Use:     "extensions",
	Aliases: []string{"ext"},
	Short:   constant.App + " extensions",
}

func init() {
	extensionsCmd.AddCommand(extensionsNewCmd)

	extensionsNewCmd.Flags().BoolP("print-path", "p", false, "print path")
}

var extensionsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new extension",
	Run: func(cmd *cobra.Command, args []string) {
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
		handleErr(err)

		lang, _ := language.FromName(answers.Language)

		var domain = passport.Domain(answers.Domain)

		p := passport.Passport{
			Name:     answers.Name,
			ID:       passport.ToID(answers.Name),
			About:    answers.About,
			Version:  *semver.MustParse("0.1.0"),
			Domain:   passport.Domain(answers.Domain),
			Language: *lang,
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
		handleErr(filesystem.Api().MkdirAll(path, os.ModePerm))

		data, err := toml.Marshal(p)
		handleErr(err)

		handleErr(filesystem.Api().WriteFile(filepath.Join(path, constant.Passport), data, os.ModePerm))

		var scraperTemplate string

		switch domain {
		case passport.DomainAnime:
			scraperTemplate = anime.Template
		case passport.DomainMovies:
			scraperTemplate = movie.Template
		default:
			handleErr(errors.New("domain not supported"))
		}

		handleErr(filesystem.Api().WriteFile(filepath.Join(path, constant.Scraper), []byte(scraperTemplate), os.ModePerm))

		fmt.Printf(
			"%s Created %s extension for %s domain\n",
			style.Fg(color.Green)(icon.Check),
			style.Fg(color.Purple)(answers.Name),
			style.New().Foreground(color.Yellow).Bold(true).Render(answers.Domain),
		)

		if printPath := lo.Must(cmd.Flags().GetBool("print-path")); printPath {
			fmt.Println(path)
		}
	},
}
