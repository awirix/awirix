package cmd

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/vivi-app/vivi/color"
	"github.com/vivi-app/vivi/extensions/extension"
	"github.com/vivi-app/vivi/extensions/manager"
	"github.com/vivi-app/vivi/filesystem"
	"github.com/vivi-app/vivi/icon"
	"github.com/vivi-app/vivi/style"
)

func completionExtensionsIDs(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	exts, err := manager.InstalledExtensions()
	if err != nil {
		exts = []*extension.Extension{}
	}

	return lo.Map(exts, func(e *extension.Extension, _ int) string {
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
	case pathFlag != nil && pathFlag.Changed:
		path := pathFlag.Value.String()
		isDir, err := filesystem.Api().IsDir(path)
		handleErr(err)

		if !isDir {
			handleErr(fmt.Errorf("path %s is not a directory", path))
		}

		ext, err = extension.New(path)
		handleErr(err)
	case idFlag != nil && idFlag.Changed:
		var err error
		id := idFlag.Value.String()
		ext, err = manager.GetExtensionByID(id)

		handleErr(err)
	}

	return ext
}

func printWarning(message string) {
	fmt.Printf("%s %s\n", style.Fg(color.Yellow)(icon.Warn), message)
}

func printError(message string) {
	fmt.Printf("%s %s\n", style.Fg(color.Red)(icon.Cross), message)
}

func printSuccess(message string) {
	fmt.Printf("%s %s\n", style.Fg(color.Green)(icon.Check), message)
}
