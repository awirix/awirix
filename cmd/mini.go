package cmd

import (
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/mini"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(miniCmd)
	miniCmd.Flags().BoolP("debug", "d", false, "enable debug mode")
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Run a mini version of " + app.Name,
	Run: func(cmd *cobra.Command, args []string) {
		options := &mini.Options{
			Debug: lo.Must(cmd.Flags().GetBool("debug")),
		}

		err := mini.Run(options)
		handleErr(err)
	},
}
