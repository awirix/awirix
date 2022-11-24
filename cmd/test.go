package cmd

import (
	"fmt"
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
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]

		isDir, err := filesystem.Api().IsDir(path)
		handleErr(err)

		if lo.Must(cmd.Flags().GetBool("passport")) {
			p, err := passport.FromPath(path)
			handleErr(err)
			fmt.Println(p.String())
		} else {
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
