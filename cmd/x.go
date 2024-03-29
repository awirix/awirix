package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/extensions/extension"
	"github.com/awirix/awirix/extensions/manager"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/lualib"
	"github.com/awirix/awirix/style"
	"github.com/awirix/awirix/text"
	"github.com/awirix/lua"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"regexp"
)

func init() {
	rootCmd.AddCommand(xCmd)
}

var xCmd = &cobra.Command{
	Use:   "x",
	Short: "Manage extensions",
	Args:  cobra.NoArgs,
}

func init() {
	xCmd.AddCommand(xNewCmd)

	xNewCmd.Flags().BoolP("print-path", "p", false, "print path")
}

var xNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new extension",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		ext, err := extension.GenerateInteractive()
		handleErr(err)

		fmt.Printf(
			"%s Created %s extension\n",
			style.Fg(color.Green)(icon.Check),
			style.Fg(color.Purple)(ext.String()),
		)

		if printPath := lo.Must(cmd.Flags().GetBool("print-path")); printPath {
			fmt.Println(ext.Path())
		}
	},
}

func init() {
	xCmd.AddCommand(xLsCmd)
}

var xLsCmd = &cobra.Command{
	Use:   "ls",
	Short: "List installed extensions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.Installed()
		handleErr(err)

		for _, e := range extensions {
			fmt.Printf(
				"%s %s\n",
				style.Fg(color.Purple)(e.Passport().Name),
				style.Bold(e.Passport().Version.String()),
			)
		}

		fmt.Printf(
			"%s %s installed\n",
			style.Fg(color.Pink)(icon.Heart),
			text.Quantify(len(extensions), "extension", "extensions"),
		)
	},
}

func init() {
	xCmd.AddCommand(xRemoveCmd)

	xRemoveCmd.Flags().StringP("id", "i", "", "id of the extension to remove")
}

var xRemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove an extension",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.Installed()
		handleErr(err)

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

			handleErr(manager.Remove(toRemove))

			return
		}

		var nameExtensionMap = make(map[string]*extension.Extension)

		for _, ext := range extensions {
			nameExtensionMap[ext.String()] = ext
		}

		var selected []string
		err = survey.AskOne(&survey.MultiSelect{
			Message: "Select extensions to remove",
			Options: lo.Keys(nameExtensionMap),
		}, &selected)
		handleErr(err)

		var confirm bool
		err = survey.AskOne(&survey.Confirm{
			Message: fmt.Sprintf("Remove %s?", text.Quantify(len(selected), "extension", "extensions")),
			Default: false,
		}, &confirm)
		handleErr(err)

		if !confirm {
			fmt.Printf("%s OK, aborted\n", style.Fg(color.Green)(icon.Check))
			return
		}

		for _, s := range selected {
			err = manager.Remove(nameExtensionMap[s])
			handleErr(err)
		}

		fmt.Printf(
			"%s Successfully removed %s\n",
			style.Fg(color.Green)(icon.Check),
			text.Quantify(len(selected), "extension", "extensions"),
		)
	},
}

func init() {
	xCmd.AddCommand(xSelCmd)

	xSelCmd.Flags().StringP("path", "p", "", "path to the extension")
	xSelCmd.Flags().StringP("id", "i", "", "id of the extension")
	xSelCmd.MarkFlagsMutuallyExclusive("path", "id")
	xSelCmd.MarkFlagDirname("path")
	xSelCmd.RegisterFlagCompletionFunc("id", completionExtensionsIDs)

	xSelCmd.Flags().Bool("run", false, "run the selected extension")
	xSelCmd.Flags().Bool("test", false, "test the selected extension")
	xSelCmd.Flags().Bool("info", false, "show info about the extension")
	xSelCmd.Flags().BoolP("json", "j", false, "output in json format")

	xSelCmd.MarkFlagsMutuallyExclusive("run", "test", "info")
}

var xSelCmd = &cobra.Command{
	Use:   "sel",
	Short: "Select an extension to perform an action on",
	Args:  cobra.NoArgs,
	PreRunE: preRunERequiredMutuallyExclusiveFlags(
		[]string{"path", "id"},
		[]string{"run", "test", "info"},
	),
	Run: func(cmd *cobra.Command, args []string) {
		ext := loadExtension(cmd.Flag("path"), cmd.Flag("id"))

		switch {
		case lo.Must(cmd.Flags().GetBool("run")):
			handleErr(ext.LoadCore(true))

		case lo.Must(cmd.Flags().GetBool("test")):
			handleErr(ext.LoadTester(true))
			handleErr(ext.Tester().Test())

		case lo.Must(cmd.Flags().GetBool("info")):
			fmt.Println(ext.Passport().Info())
		}
	},
}

func init() {
	xCmd.AddCommand(xAddCmd)
}

var xAddCmd = &cobra.Command{
	Use:     "add",
	Short:   "Install an extension",
	Aliases: []string{"install"},
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		arg := args[0]

		var url string

		if text.IsURLStrict(arg) {
			url = arg
		} else if matched, _ := regexp.MatchString(`^[\w-]+/[\w-]+$`, arg); matched {
			url = fmt.Sprintf("https://github.com/%s", arg)
		}

		ext, err := manager.Add(url, &manager.AddOptions{})
		handleErr(err)

		fmt.Printf("%s Successfully installed %s\n", style.Fg(color.Green)(icon.Check), style.Fg(color.Purple)(ext.String()))
	},
}

func init() {
	xCmd.AddCommand(xUpCmd)
	xUpCmd.Flags().BoolP("verbose", "v", false, "print skipped extensions")
}

var xUpCmd = &cobra.Command{
	Use:   "up",
	Short: "Update extensions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		extensions, err := manager.Installed()
		handleErr(err)

		for _, ext := range extensions {
			if ext.Passport().Repository.IsAbsent() {
				if lo.Must(cmd.Flags().GetBool("verbose")) {
					printWarning(fmt.Sprintf("skipping %s, no repository specified", style.Fg(color.Purple)(ext.String())))
				}
				continue
			}

			updated, err := manager.Update(ext)

			if err != nil {
				printError(fmt.Sprintf("failed to update %s: %s", style.Fg(color.Purple)(ext.String()), err))
				continue
			}

			var outcome string
			if updated.Passport().Version.Compare(&ext.Passport().Version) > 0 {
				outcome = fmt.Sprintf("updated %s => %s", ext.Passport().Version.String(), updated.Passport().Version.String())
			} else {
				outcome = "already up to date"
			}

			printSuccess(fmt.Sprintf("%s %s", style.Fg(color.Purple)(ext.String()), outcome))
		}
	},
}

func init() {
	xCmd.AddCommand(xHealthCmd)
}

var xHealthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check the health of the extensions",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		manager.Health(cmd.OutOrStdout())
	},
}

func init() {
	xLibCmd.Flags().BoolP("stdout", "s", false, "write output to stdout")
	xCmd.AddCommand(xLibCmd)
}

var xLibCmd = &cobra.Command{
	Use:   "lib",
	Short: "Generate library hints",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		state := lua.NewState(nil)
		lib := lualib.Lib(state)
		doc := lib.LuaDoc()

		if lo.Must(cmd.Flags().GetBool("stdout")) {
			fmt.Println(doc)
		} else {
			filename := app.Name + ".lua"
			err := filesystem.Api().WriteFile(filename, []byte(doc), os.ModePerm)
			handleErr(err)

			_, err = log.WriteSuccessf(os.Stderr, "%s generated\n", filename)
			handleErr(err)
		}
	},
}
