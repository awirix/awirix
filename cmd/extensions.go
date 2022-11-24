package cmd

import (
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
	"github.com/vivi-app/vivi/scraper"
	"github.com/vivi-app/vivi/scraper/test"
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
		handleErr(filesystem.Api().MkdirAll(path, os.ModePerm))

		data, err := toml.Marshal(p)
		handleErr(err)

		handleErr(filesystem.Api().WriteFile(filepath.Join(path, constant.Passport), data, os.ModePerm))

		scraperTemplate, err := scraper.GenerateTemplate(domain)
		handleErr(err)

		// TODO: dry
		handleErr(filesystem.Api().WriteFile(filepath.Join(path, constant.Scraper), []byte(scraperTemplate), os.ModePerm))
		handleErr(filesystem.Api().WriteFile(filepath.Join(path, constant.Test), []byte(test.Template), os.ModePerm))

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

func init() {
	extensionsCmd.AddCommand(extensionsListCmd)
}

var extensionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed extensions",
	Run: func(cmd *cobra.Command, args []string) {
		// get all dirs
		descriptors, err := filesystem.Api().ReadDir(where.Extensions())
		handleErr(err)

		var (
			passports = make(map[passport.Domain][]*passport.Passport)
			count     int
		)

		for _, d := range descriptors {
			if d.IsDir() {
				if p, err := passport.FromPath(filepath.Join(where.Extensions(), d.Name())); err == nil {
					passports[p.Domain] = append(passports[p.Domain], p)
					count++
				}
			}
		}

		if len(passports) == 0 {
			fmt.Println("No extensions installed")
			return
		}

		printForDomain := func(d passport.Domain) {
			fmt.Println(
				style.
					New().
					Foreground(color.Yellow).
					Bold(true).
					Render(util.Capitalize(string(d))),
			)

			for _, p := range passports[d] {
				fmt.Printf(
					"%s %s %s\n",
					style.Fg(color.Purple)(p.Name),
					style.Bold(p.Version.String()),
					style.Faint(p.About),
				)
			}
		}

		for _, domain := range passport.Domains {
			if _, ok := passports[domain]; ok {
				printForDomain(domain)
				fmt.Println()
			}
		}

		fmt.Printf(
			"%s %s installed\n",
			style.Fg(color.Pink)(icon.Heart),
			util.Quantify(count, "extension", "extensions"),
		)
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsRemoveCmd)

	extensionsRemoveCmd.Flags().StringP("name", "n", "", "Name of the extension to remove")
}

var extensionsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove an extension",
	Run: func(cmd *cobra.Command, args []string) {
		if name := lo.Must(cmd.Flags().GetString("name")); name != "" {
			path := filepath.Join(where.Extensions(), util.SanitizeFilename(name))
			if err := filesystem.Api().Remove(filepath.Join(where.Extensions(), path)); err != nil {
				handleErr(err)
			}

			return
		}

		// get all dirs
		descriptors, err := filesystem.Api().ReadDir(where.Extensions())
		handleErr(err)

		var (
			extensions = make(map[string]string)
		)

		for _, d := range descriptors {
			if !d.IsDir() {
				continue
			}

			path := filepath.Join(where.Extensions(), d.Name())
			p, err := passport.FromPath(path)
			if err != nil {
				continue
			}

			extensions[p.Name] = path
		}

		var selected []string

		err = survey.AskOne(&survey.MultiSelect{
			Message: "Select extensions to remove",
			Options: lo.Keys(extensions),
		}, &selected)

		handleErr(err)

		for _, s := range selected {
			if err := filesystem.Api().RemoveAll(extensions[s]); err != nil {
				handleErr(err)
			}
		}

		fmt.Printf(
			"%s Successfully removed %s\n",
			style.Fg(color.Green)(icon.Check),
			util.Quantify(len(selected), "extension", "extensions"),
		)
	},
}
