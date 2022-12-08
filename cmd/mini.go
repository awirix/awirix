package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/mini"
)

func init() {
	rootCmd.AddCommand(miniCmd)
}

var miniCmd = &cobra.Command{
	Use:   "mini",
	Short: "Run a mini version of Vivi",
	Long:  `Run a mini version of Vivi`,
	Run: func(cmd *cobra.Command, args []string) {
		err := mini.Run(nil)
		handleErr(err)
	},
}
