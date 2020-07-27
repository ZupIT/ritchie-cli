package cmd

import (
	"github.com/spf13/cobra"
)

// NewSetCmd creates new cmd instance
func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set SUBCOMMAND",
		Short:   "Set contexts, credentials and priorities",
		Long:    "Set contexts, credentials and priorities for formula repositories.",
		Example: "rit set context",
	}
}
