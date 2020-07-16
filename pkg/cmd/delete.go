package cmd

import "github.com/spf13/cobra"

// NewDeleteCmd create a new delete instance
func NewDeleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "delete SUBCOMMAND",
		Short: "Delete objects (contexts, repositories)",
		Long:  `Delete objects like contexts and repo`,
	}
}
