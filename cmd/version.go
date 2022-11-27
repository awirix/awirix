package cmd

import (
	"github.com/samber/lo"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/style"
	"html/template"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/constant"
)

func init() {
	rootCmd.AddCommand(versionCmd)

	versionCmd.Flags().BoolP("short", "s", false, "print the version number only")
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of the " + constant.App,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if lo.Must(cmd.Flags().GetBool("short")) {
			_, err := cmd.OutOrStdout().Write([]byte(constant.Version + "\n"))
			handleErr(err)
			return
		}

		versionInfo := struct {
			Version  string
			OS       string
			Arch     string
			App      string
			Compiler string
		}{
			Version:  constant.Version,
			App:      constant.App,
			OS:       runtime.GOOS,
			Arch:     runtime.GOARCH,
			Compiler: runtime.Compiler,
		}

		t, err := template.New("version").Funcs(map[string]any{
			"faint":   style.Faint,
			"bold":    style.Bold,
			"magenta": style.Fg(color.Purple),
		}).Parse(`{{ magenta "▇▇▇" }} {{ magenta .App }} 

  {{ faint "Version" }}  {{ bold .Version }}
  {{ faint "Platform" }} {{ bold .OS }}/{{ bold .Arch }}
  {{ faint "Compiler" }} {{ bold .Compiler }}
`)
		handleErr(err)
		handleErr(t.Execute(cmd.OutOrStdout(), versionInfo))
	},
}
