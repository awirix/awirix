package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/constant"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/where"
	"os"
)

type whereTarget struct {
	name  string
	where func() string

	argShort, argLong string
}

// Specify what paths to show
var wherePaths = []whereTarget{
	{"Downloads", where.Downloads, "d", "downloads"},
	{"Extensions", where.Extensions, "e", "extensions"},
	{"Config", where.Config, "c", "config"},
	{"Logs", where.Logs, "l", "logs"},
}

func init() {
	rootCmd.AddCommand(whereCmd)

	for _, n := range wherePaths {
		if n.argShort != "" {
			whereCmd.Flags().BoolP(n.argLong, n.argShort, false, n.name+" path")
		} else {
			whereCmd.Flags().Bool(n.argLong, false, n.name+" path")
		}
	}

	whereCmd.MarkFlagsMutuallyExclusive(lo.Map(wherePaths, func(t whereTarget, _ int) string {
		return t.argLong
	})...)

	whereCmd.SetOut(os.Stdout)
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show the paths for a files related to the " + constant.App,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		headerStyle := style.NewColored(color.HiPurple, "").Bold(true).Render
		argStyle := style.Fg(color.Yellow)

		for _, n := range wherePaths {
			if lo.Must(cmd.Flags().GetBool(n.argLong)) {
				cmd.Println(n.where())
				return
			}
		}

		for i, n := range wherePaths {
			cmd.Printf("%s %s\n", headerStyle(n.name+"?"), argStyle("--"+n.argLong))
			cmd.Println(n.where())

			if i < len(wherePaths)-1 {
				cmd.Println()
			}
		}
	},
}
