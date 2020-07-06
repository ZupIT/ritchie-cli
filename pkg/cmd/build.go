package cmd

import "github.com/spf13/cobra"

func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build SUB_COMMAND",
		Short: "Build objects (formula and others)",
		Long:  `This is a root command, to use it, add a sub command. For example, rit build formula.`,
	}
}
