package cmd

import "github.com/spf13/cobra"

func NewShowCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "show SUB_COMMAND",
		Short:   "Show context",
		Long:    "Show current context",
		Example: "rit show context",
	}
}
