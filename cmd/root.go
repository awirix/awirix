package cmd

import (
	"fmt"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/app"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/log"
	"github.com/vivi-app/vivi/style"
	"github.com/vivi-app/vivi/where"
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
		log.Error(err)
		_, _ = fmt.Fprintf(
			os.Stderr,
			"%s %s\n",
			style.Fg(color.Red)(icon.Cross),
			strings.Trim(err.Error(), " \n"),
		)
		os.Exit(1)
	}
}
