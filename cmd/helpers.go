package cmd

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/vivi-app/vivi/extension"
)

func completionExtensionsNames(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	return lo.Map(extension.ListInstalled(), func(e *extension.Extension, _ int) string {
		return e.String()
	}), cobra.ShellCompDirectiveNoFileComp
}

func preRunERequiredMutuallyExclusiveFlags(flags ...string) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		for _, flag := range flags {
			if cmd.Flag(flag).Changed {
				return nil
			}
		}

		return fmt.Errorf("one of the flags %v is required", flags)
	}
}
