package cmd

import "github.com/spf13/cobra"

func NewShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show SUB_COMMAND",
		Short: "Show contexts",
		Long:  `Show objects like context, organization etc.`,
	}
}
