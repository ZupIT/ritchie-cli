package cmd

import (
	"github.com/spf13/cobra"
)

// NewCreateCmd creates new cmd instance
func NewCreateCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "create SUBCOMMAND",
		Short:   "Create formulas",
		Long:    "Create a new formula.",
		Example: "rit create formula",
	}
}
