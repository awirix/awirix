package cmd

import (
	"fmt"
	"github.com/pelletier/go-toml/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/lualib"
	"github.com/vivi-app/vivi/passport"
	lua "github.com/yuin/gopher-lua"
	"path/filepath"
)

func init() {
	rootCmd.AddCommand(testCmd)

	testCmd.Flags().BoolP("passport", "p", false, "test passport")
	testCmd.Flags().BoolP("run", "r", false, "run scraper")

	testCmd.MarkFlagsMutuallyExclusive("passport", "run")
}

var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test scraper",
	Args: func(cmd *cobra.Command, args []string) error {
		// expect 1 argument
		if len(args) != 1 {
			return fmt.Errorf("accepts 1 arg(s), received %d", len(args))
		}

		// check if path exists
		path := args[0]
		exists, err := filesystem.Api().Exists(path)
		handleErr(err)
		if !exists {
			return fmt.Errorf("path does not exist: %s", path)
		}

		return nil
	},
	PreRun: func(cmd *cobra.Command, args []string) {
		if !cmd.Flags().Changed("passport") && !cmd.Flags().Changed("run") {
			handleErr(fmt.Errorf("either --passport or --run flag must be set"))
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		isDir, err := filesystem.Api().IsDir(path)
		handleErr(err)

		if lo.Must(cmd.Flags().GetBool("passport")) {
			if isDir {
				path = filepath.Join(path, constant.Passport)
			}

			data, err := filesystem.Api().ReadFile(path)
			handleErr(err)

			var passport passport.Passport
			err = toml.Unmarshal(data, &passport)
			handleErr(err)

			handleErr(passport.Validate())

			fmt.Println(passport.String())
		} else if lo.Must(cmd.Flags().GetBool("run")) {
			if isDir {
				path = filepath.Join(path, constant.Scraper)
			}

			L := lua.NewState()
			defer L.Close()

			lualib.Preload(L)
			handleErr(L.DoFile(path))
		}
	},
}
