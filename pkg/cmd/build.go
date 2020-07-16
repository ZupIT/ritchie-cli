package cmd

import "github.com/spf13/cobra"

// NewBuildCmd creates new cmd instance of build command
func NewBuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build SUB_COMMAND",
		Short: "Build formulas",
		Long:  `This is a root command, to use it, add a sub command. For example, rit build formula.`,
	}
}
