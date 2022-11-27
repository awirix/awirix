package cmd

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extension"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/style"
)

func completionExtensionsIDs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return lo.Map(extension.ListInstalled(), func(e *extension.Extension, _ int) string {
		return e.Passport().ID
	}), cobra.ShellCompDirectiveNoFileComp
}

func preRunERequiredMutuallyExclusiveFlags(flagsGroups ...[]string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		checkGroup := func(group []string) error {
			for _, flag := range group {
				if cmd.Flag(flag).Changed {
					return nil
				}
			}

			return fmt.Errorf("one of the flags %v is required", group)
		}

		for _, group := range flagsGroups {
			err := checkGroup(group)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func loadExtension(pathFlag, idFlag *pflag.Flag) *extension.Extension {
	var ext *extension.Extension

	switch {
	case pathFlag.Changed:
		path := pathFlag.Value.String()
		isDir, err := filesystem.Api().IsDir(path)
		handleErr(err)

		if !isDir {
			handleErr(fmt.Errorf("path %s is not a directory", path))
		}

		ext, err = extension.NewFromPath(path)
		handleErr(err)
	case idFlag.Changed:
		var ok bool
		id := idFlag.Value.String()
		ext, ok = extension.NewFromID(id)

		if !ok {
			handleErr(fmt.Errorf("extension %s not found", style.Fg(color.Purple)(id)))
		}
	}

	return ext
}
