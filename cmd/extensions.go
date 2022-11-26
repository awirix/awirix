package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/extension"
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
	Args:    cobra.NoArgs,
}

func init() {
	extensionsCmd.AddCommand(extensionsNewCmd)

	extensionsNewCmd.Flags().BoolP("print-path", "p", false, "print path")
}

var extensionsNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new extension",
	Args:  cobra.NoArgs,
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
	Args:  cobra.NoArgs,
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

	extensionsRemoveCmd.Flags().StringP("id", "i", "", "id of the extension to remove")
}

var extensionsRemoveCmd = &cobra.Command{
	Use:     "remove",
	Short:   "Remove an extension",
	Aliases: []string{"rm"},
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions := extension.ListInstalled()

		if id := lo.Must(cmd.Flags().GetString("id")); id != "" {
			toRemove, ok := lo.Find(extensions, func(e *extension.Extension) bool {
				return e.String() == id
			})

			if !ok {
				handleErr(fmt.Errorf(
					"extension %s not found",
					style.Fg(color.Purple)(id),
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

	extensionsRunCmd.Flags().StringP("path", "p", "", "path to the extension to run")
	extensionsRunCmd.Flags().StringP("id", "i", "", "id of the extension to run")

	extensionsRunCmd.MarkFlagsMutuallyExclusive("path", "id")
	extensionsRunCmd.MarkFlagDirname("path")

	extensionsRunCmd.RegisterFlagCompletionFunc("id", completionExtensionsIDs)
}

var extensionsRunCmd = &cobra.Command{
	Use:     "run",
	Short:   "Run an extension for the testing purposes",
	Args:    cobra.NoArgs,
	PreRunE: preRunERequiredMutuallyExclusiveFlags("path", "id"),
	Run: func(cmd *cobra.Command, args []string) {
		ext := loadExtension(cmd.Flag("path"), cmd.Flag("id"))

		// Script will be executed when it's loaded
		// I don't like this behavior, but it cannot be avoided, so let's use it as it is
		handleErr(ext.LoadScraper())
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsInfoCmd)

	extensionsInfoCmd.Flags().StringP("path", "p", "", "Path to the extension to run")
	extensionsInfoCmd.Flags().StringP("id", "i", "", "id of the extension to run")

	extensionsInfoCmd.MarkFlagsMutuallyExclusive("path", "id")
	extensionsInfoCmd.MarkFlagDirname("path")

	extensionsInfoCmd.RegisterFlagCompletionFunc("id", completionExtensionsIDs)
}

var extensionsInfoCmd = &cobra.Command{
	Use:     "info",
	Short:   "Show extension info",
	Args:    cobra.NoArgs,
	PreRunE: preRunERequiredMutuallyExclusiveFlags("path", "id"),
	Run: func(cmd *cobra.Command, args []string) {
		ext := loadExtension(cmd.Flag("path"), cmd.Flag("id"))
		fmt.Println(ext.Passport().Info())
	},
}

func init() {
	extensionsCmd.AddCommand(extensionsTestCmd)

	extensionsTestCmd.Flags().StringP("path", "p", "", "Path to the extension to run")
	extensionsTestCmd.Flags().StringP("id", "i", "", "id of the extension to run")

	extensionsTestCmd.MarkFlagsMutuallyExclusive("path", "id")
	extensionsTestCmd.MarkFlagDirname("path")

	extensionsTestCmd.RegisterFlagCompletionFunc("id", completionExtensionsIDs)
}

var extensionsTestCmd = &cobra.Command{
	Use:     "test",
	Short:   "Test an extension",
	PreRunE: preRunERequiredMutuallyExclusiveFlags("path", "id"),
	Run: func(cmd *cobra.Command, args []string) {
		ext := loadExtension(cmd.Flag("path"), cmd.Flag("id"))
		handleErr(ext.LoadTester())
		handleErr(ext.Tester().Test())
	},
}
