package cmd

import (
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/tui"
	"github.com/awirix/awirix/where"
	"os"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     strings.ToLower(app.Name),
	Short:   "Multimedia Metascraper",
	Long:    app.AsciiArt + "\nWatch anime, movies and TV shows from any source in one place.",
	Version: app.Version,
	Args:    cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		err := tui.Run(&tui.Options{AltScreen: true})
		handleErr(err)
	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:       rootCmd,
		Headings:      cc.HiCyan + cc.Bold + cc.Underline,
		Commands:      cc.HiYellow + cc.Bold,
		Example:       cc.Italic,
		ExecName:      cc.Bold,
		Flags:         cc.Bold,
		FlagsDataType: cc.Italic + cc.HiBlue,
	})

	// Clears temp files on each run.
	// It should not affect startup time since it's being run as goroutine.
	go func() {
		_ = filesystem.Api().RemoveAll(where.Temp())
	}()

	_ = rootCmd.Execute()
}

func handleErr(err error) {
	if err != nil {
		msg := strings.TrimSpace(err.Error())
		msg = strings.Trim(msg, "\n")

		log.Error(msg)
		_, _ = log.WriteErrorf(os.Stderr, msg)
		_, _ = os.Stderr.Write([]byte("\n"))

		os.Exit(1)
	}
}
