package cmd

import (
	"github.com/spf13/cobra"
)

// NewSetCmd creates new cmd instance
func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "set SUBCOMMAND",
		Short: "Set objects (context, credential)",
		Long:  `Set objects like credentials, etc.`,
	}
}
