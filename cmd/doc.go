package cmd

import (
	"fmt"
	"github.com/awirix/awirix/lualib"
	"github.com/awirix/lua"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(docCmd)
}

var docCmd = &cobra.Command{
	Use:   "doc",
	Short: "Generate documentation",
	Long:  `Generate documentation`,
}

func init() {
	docCmd.AddCommand(docLuaCmd)
}

var docLuaCmd = &cobra.Command{
	Use:   "lua",
	Short: "Generate Lua documentation that can used by language server",
	Run: func(cmd *cobra.Command, args []string) {
		state := lua.NewState(nil)
		lib := lualib.Lib(state)
		fmt.Println(lib.LuaDoc())
	},
}
