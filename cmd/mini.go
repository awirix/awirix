package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/mini"
)

func init() {
	rootCmd.AddCommand(miniCmd)
	miniCmd.Flags().BoolP("debug", "d", false, "enable debug mode")
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Run a mini version of Vivi",
	Long:  `Run a mini version of Vivi`,
	Run: func(cmd *cobra.Command, args []string) {
		options := &mini.Options{
			Debug: lo.Must(cmd.Flags().GetBool("debug")),
		}

		err := mini.Run(options)
		handleErr(err)
	},
}
