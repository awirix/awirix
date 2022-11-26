package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/extension"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/passport"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/util"
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
		ext, err := extension.NewInteractive()
		handleErr(err)

		fmt.Printf(
			"%s Created %s extension for %s domain\n",
			style.Fg(color.Green)(icon.Check),
			style.Fg(color.Purple)(ext.String()),
			style.New().Foreground(color.Yellow).Bold(true).Render(string(ext.Passport().Domain)),
		)

		if printPath := lo.Must(cmd.Flags().GetBool("print-path")); printPath {
			fmt.Println(ext.Path())
		}
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsListCmd)
}

var extensionsListCmd = &cobra.Command{
	Use:   "list",
	Short: "ListInstalled installed extensions",
	Run: func(cmd *cobra.Command, args []string) {
		extensions := extension.ListInstalled()

		var byDomain = make(map[passport.Domain][]*extension.Extension)

		for _, ext := range extensions {
			byDomain[ext.Passport().Domain] = append(byDomain[ext.Passport().Domain], ext)
		}

		printForDomain := func(d passport.Domain) {
			fmt.Println(
				style.
					New().
					Foreground(color.Yellow).
					Bold(true).
					Render(util.Capitalize(string(d))),
			)

			for _, e := range byDomain[d] {
				fmt.Printf(
					"%s %s %s\n",
					style.Fg(color.Purple)(e.String()),
					style.Bold(e.Passport().Version.String()),
					style.Faint(e.Passport().About),
				)
			}
		}

		for _, domain := range passport.Domains {
			if _, ok := byDomain[domain]; ok {
				printForDomain(domain)
				fmt.Println()
			}
		}

		fmt.Printf(
			"%s %s installed\n",
			style.Fg(color.Pink)(icon.Heart),
			util.Quantify(len(extensions), "extension", "extensions"),
		)
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsRemoveCmd)

	extensionsRemoveCmd.Flags().StringP("name", "n", "", "Name of the extension to remove")
}

var extensionsRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an extension",
	Aliases: []string{"rm"},
	Run: func(cmd *cobra.Command, args []string) {
		extensions := extension.ListInstalled()

		if name := lo.Must(cmd.Flags().GetString("name")); name != "" {
			toRemove, ok := lo.Find(extensions, func(e *extension.Extension) bool {
				return e.String() == name
			})

			if !ok {
				handleErr(fmt.Errorf(
					"extension %s not found",
					style.Fg(color.Purple)(name),
				))
			}

			handleErr(toRemove.Uninstall())

			return
		}

		var nameExtensionMap = make(map[string]*extension.Extension)

		for _, ext := range extensions {
			nameExtensionMap[ext.String()] = ext
		}

		var selected []string
		err := survey.AskOne(&survey.MultiSelect{
			Message: "Select extensions to remove",
			Options: lo.Keys(nameExtensionMap),
		}, &selected)
		handleErr(err)

		var confirm bool
		err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Remove %s?", strings.Join(lo.Map(selected, func(s string, _ int) string {
				return style.Fg(color.Purple)(s)
			}), ", ")),
			Default: false,
		}, &confirm)
		handleErr(err)

		if !confirm {
			fmt.Printf("%s OK, aborted\n", style.Fg(color.Green)(icon.Check))
			return
		}

		for _, s := range selected {
			handleErr(nameExtensionMap[s].Uninstall())
		}

		fmt.Printf(
			"%s Successfully removed %s\n",
			style.Fg(color.Green)(icon.Check),
			util.Quantify(len(selected), "extension", "extensions"),
		)
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsRunCmd)

	extensionsRunCmd.Flags().StringP("path", "p", "", "Path to the extension to run")
	extensionsRunCmd.Flags().StringP("name", "n", "", "Name of the extension to run")

	extensionsRunCmd.MarkFlagsMutuallyExclusive("path", "name")
	extensionsRunCmd.MarkFlagDirname("path")

	extensionsRunCmd.RegisterFlagCompletionFunc("name", completionExtensionsNames)
}

var extensionsRunCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run an extension for the testing purposes",
	PreRunE: preRunERequiredMutuallyExclusiveFlags("path", "name"),
	Run: func(cmd *cobra.Command, args []string) {
		var (
			ext *extension.Extension
			err error
		)

		if cmd.Flags().Changed("path") {
			path := lo.Must(cmd.Flags().GetString("path"))
			isDir, err := filesystem.Api().IsDir(path)
			handleErr(err)

			if !isDir {
				handleErr(fmt.Errorf("path %s is not a directory", path))
				return
			}

			ext, err = extension.FromPath(path)
			handleErr(err)
		} else {
			var ok bool
			name := lo.Must(cmd.Flags().GetString("name"))
			ext, ok = extension.FromLocalName(name)

			if !ok {
				handleErr(fmt.Errorf("extension %s not found", style.Fg(color.Purple)(name)))
			}
		}

		// Script will be executed when it's loaded
		// I don't like this behavior, but it cannot be avoided, so let's use it as it is
		err = ext.LoadScraper()
		handleErr(err)
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsInfoCmd)

	extensionsInfoCmd.Flags().StringP("path", "p", "", "Path to the extension to run")
	extensionsInfoCmd.Flags().StringP("name", "n", "", "Name of the extension to run")

	extensionsInfoCmd.MarkFlagsMutuallyExclusive("path", "name")
	extensionsInfoCmd.MarkFlagDirname("path")

	extensionsInfoCmd.RegisterFlagCompletionFunc("name", completionExtensionsNames)
}

var extensionsInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Show extension info",
	PreRunE: preRunERequiredMutuallyExclusiveFlags("path", "name"),
	Run: func(cmd *cobra.Command, args []string) {
		var ext *extension.Extension

		switch {
		case cmd.Flags().Changed("path"):
			path := lo.Must(cmd.Flags().GetString("path"))
			isDir, err := filesystem.Api().IsDir(path)
			handleErr(err)

			if !isDir {
				handleErr(fmt.Errorf("path %s is not a directory", path))
				return
			}

			ext, err = extension.FromPath(path)
			handleErr(err)
		case cmd.Flags().Changed("name"):
			var ok bool

			name := lo.Must(cmd.Flags().GetString("name"))
			ext, ok = extension.FromLocalName(name)

			if !ok {
				handleErr(fmt.Errorf("extension %s not found", style.Fg(color.Purple)(name)))
			}
		}

		fmt.Println(ext.Passport().Info())
	},
}
