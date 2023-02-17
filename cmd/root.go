package cmd

import (
	"github.com/awirix/awirix/app"
	"github.com/awirix/awirix/filesystem"
	"github.com/awirix/awirix/log"
	"github.com/awirix/awirix/tui"
	"github.com/awirix/awirix/where"
	cc "github.com/ivanpirog/coloredcobra"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func init() {
	rootCmd.Flags().BoolP("version", "v", false, "version for "+app.Name)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   strings.ToLower(app.Name),
	Short: "Multimedia Metascraper",
	Long:  app.AsciiArt + "\nWatch anime, movies and TV shows from any source in one place.",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if version, _ := cmd.Flags().GetBool("version"); version {
			versionCmd.Run(versionCmd, args)
			return
		}

		err := tui.Run(&tui.Options{AltScreen: true})
		handleErr(err)
	},
}

func Execute() {
	cc.Init(&cc.Config{
		RootCmd:         rootCmd,
		Headings:        cc.HiCyan + cc.Bold + cc.Underline,
		Commands:        cc.HiYellow + cc.Bold,
		Example:         cc.Italic,
		ExecName:        cc.Bold,
		Flags:           cc.Bold,
		FlagsDataType:   cc.Italic + cc.HiBlue,
		NoExtraNewlines: true,
		NoBottomNewline: true,
	})

	// Clears temp files on each run.
	// It should not affect startup time since it's being run as a separate goroutine.
	go func() {
		_ = filesystem.Api().RemoveAll(where.Temp())
	}()

	_ = rootCmd.Execute()
}

func handleErr(err error) {
	if err == nil {
		return
	}

	msg := strings.TrimSpace(err.Error())
	msg = strings.Trim(msg, "\n")

	log.Error(msg)
	_, _ = log.WriteErrorf(os.Stderr, msg)
	_, _ = os.Stderr.Write([]byte("\n"))

	os.Exit(1)
}
