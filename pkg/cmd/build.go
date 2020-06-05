package cmd

import "github.com/spf13/cobra"

func NewTestCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "build SUB_COMMAND",
		Short: "This is a root command, needs a sub command",
		Long:  `This is a root command, to use it, add a sub command. For example, rit build formula.`,
	}
}