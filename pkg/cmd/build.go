package cmd

import "github.com/spf13/cobra"

func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "build SUB_COMMAND",
		Short:   "Build formulas - with subcomando formulas",
		Long:    "This is a root command, to use it, add a sub command.",
		Example: "rit build formula",
	}
}
