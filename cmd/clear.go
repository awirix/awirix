package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/color"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/icon"
	"github.com/awirix/awirix/style"
	"github.com/awirix/awirix/text"
	"github.com/awirix/awirix/where"
)

type clearTarget struct {
	name  string
	clear func() error
}

// Specify what can be cleared
var clearTargets = []clearTarget{
	{"cache", func() error {
		return filesystem.Api().RemoveAll(where.Cache())
	}},
	{"logs", func() error {
		return filesystem.Api().RemoveAll(where.Logs())
	}},
}

func init() {
	rootCmd.AddCommand(clearCmd)
	for _, n := range clearTargets {
		clearCmd.Flags().BoolP(n.name, string(n.name[0]), false, "clear "+n.name)
	}
}

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clears sidelined files produced by the " + app.Name,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		successStyle := style.Fg(color.Green)
		var didSomething bool
		for _, n := range clearTargets {
			if lo.Must(cmd.Flags().GetBool(n.name)) {
				handleErr(n.clear())
				cmd.Printf("%s %s cleared\n", successStyle(icon.Check), text.Capitalize(n.name))
				didSomething = true
			}
		}

		if !didSomething {
			_ = cmd.Help()
		}
	},
}
