package cmd

import (
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/extensions/passport"
	"github.com/vivi-app/vivi/schema"
)

func init() {
	rootCmd.AddCommand(schemaCmd)
}

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "JSON Schema for different data",
}

func init() {
	schemaCmd.AddCommand(schemaPassportCmd)
}

var schemaPassportCmd = &cobra.Command{
	Use:   "passport",
	Short: "JSON Schema for the passport.json file",
	Run: func(cmd *cobra.Command, args []string) {
		data, err := schema.Reflect[passport.Passport]()
		handleErr(err)

		_, err = cmd.OutOrStdout().Write(data)
		handleErr(err)
	},
}
