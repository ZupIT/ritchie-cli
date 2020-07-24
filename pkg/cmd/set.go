package cmd

import (
	"github.com/spf13/cobra"
)

// NewSetCmd creates new cmd instance
func NewSetCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "set SUBCOMMAND",
		Short:   "Set objects (context, credential, repo-priority)",
		Long:    "Set contexts, credentials and priorities for formula repositories.",
		Example: "rit set context",
	}
}
