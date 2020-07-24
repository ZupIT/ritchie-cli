package cmd

import "github.com/spf13/cobra"

// NewUpdateCmd create a new update instance
func NewUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "update SUBCOMMAND",
		Short:   "Update repositories",
		Long:    "Update repositories.",
		Example: "rit update repo",
	}
}
